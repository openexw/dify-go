package sse

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
)

var (
	headerID    = []byte("id:")
	headerData  = []byte("data:")
	headerEvent = []byte("event:")
	headerRetry = []byte("retry:")
)

type Stream struct {
	EncodingBase64 bool
	LastEventID    atomic.Value // []byte
	resp           *http.Response
}

func New(resp *http.Response) *Stream {
	return &Stream{
		resp: resp,
	}
}

// Subscribe to a data stream with context
func (c *Stream) Subscribe(ctx context.Context, handler func(msg *Event)) (err error) {
	if c.resp.StatusCode != http.StatusOK {
		c.resp.Body.Close()
		return fmt.Errorf("could not connect to stream: %s", http.StatusText(c.resp.StatusCode))
	}
	defer c.resp.Body.Close()

	reader := NewEventStreamReader(c.resp.Body, 1<<16)
	eventChan, errorChan := c.startReadLoop(reader)

	for {
		select {
		case err = <-errorChan:
			return err
		case msg := <-eventChan:
			handler(msg)
		}
	}
}

func (c *Stream) startReadLoop(reader *EventStreamReader) (chan *Event, chan error) {
	outCh := make(chan *Event)
	erChan := make(chan error)
	go c.readLoop(reader, outCh, erChan)
	return outCh, erChan
}

func (c *Stream) readLoop(reader *EventStreamReader, outCh chan *Event, erChan chan error) {
	for {
		// Read each new line and process the type of event
		event, err := reader.ReadEvent()
		if err != nil {
			if err == io.EOF {
				erChan <- nil
				return
			}
			erChan <- err
			return
		}

		// If we get an error, ignore it.
		var msg *Event
		if msg, err = c.processEvent(event); err == nil {
			if len(msg.ID) > 0 {
				c.LastEventID.Store(msg.ID)
			} else {
				msg.ID, _ = c.LastEventID.Load().([]byte)
			}

			// Send downstream if the event has something useful
			if msg.hasContent() {
				outCh <- msg
			}
		}
	}
}

func trimHeader(size int, data []byte) []byte {
	if data == nil || len(data) < size {
		return data
	}

	data = data[size:]
	// Remove optional leading whitespace
	if len(data) > 0 && data[0] == 32 {
		data = data[1:]
	}
	// Remove trailing new line
	if len(data) > 0 && data[len(data)-1] == 10 {
		data = data[:len(data)-1]
	}
	return data
}

func (c *Stream) processEvent(msg []byte) (event *Event, err error) {
	var e Event

	if len(msg) < 1 {
		return nil, errors.New("event message was empty")
	}

	// Normalize the crlf to lf to make it easier to split the lines.
	// Split the line by "\n" or "\r", per the spec.
	for _, line := range bytes.FieldsFunc(msg, func(r rune) bool { return r == '\n' || r == '\r' }) {
		switch {
		case bytes.HasPrefix(line, headerID):
			e.ID = append([]byte(nil), trimHeader(len(headerID), line)...)
		case bytes.HasPrefix(line, headerData):
			// The spec allows for multiple data fields per event, concatenated them with "\n".
			e.Data = append(e.Data[:], append(trimHeader(len(headerData), line), byte('\n'))...)
		// The spec says that a line that simply contains the string "data" should be treated as a data field with an empty body.
		case bytes.Equal(line, bytes.TrimSuffix(headerData, []byte(":"))):
			e.Data = append(e.Data, byte('\n'))
		case bytes.HasPrefix(line, headerEvent):
			e.Event = append([]byte(nil), trimHeader(len(headerEvent), line)...)
		case bytes.HasPrefix(line, headerRetry):
			e.Retry = append([]byte(nil), trimHeader(len(headerRetry), line)...)
		default:
			// Ignore any garbage that doesn't match what we're looking for.
		}
	}

	// Trim the last "\n" per the spec.
	e.Data = bytes.TrimSuffix(e.Data, []byte("\n"))

	if c.EncodingBase64 {
		buf := make([]byte, base64.StdEncoding.DecodedLen(len(e.Data)))

		n, err := base64.StdEncoding.Decode(buf, e.Data)
		if err != nil {
			err = fmt.Errorf("failed to decode event message: %s", err)
		}
		e.Data = buf[:n]
	}
	return &e, err
}
