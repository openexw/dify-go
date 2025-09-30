package chatflow

type Chatflow interface {
}
type chatflow struct{}

func NewChatflow() Chatflow {
	return &chatflow{}
}
