package workflow

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	workflowv1 "github.com/openexw/dify-go/api/workflow/v1"
	"golang.org/x/net/context"
	"resty.dev/v3"
)

type Workflow interface {
	// Run workflow. Cannot be executed without a published workflow.
	Run(ctx context.Context, request workflowv1.RunRequest) (*workflowv1.RunBlockingResponse, error)
	// RunStream workflow. Cannot be executed without a published workflow.
	RunStream(ctx context.Context, request workflowv1.RunRequest, fn func(v any)) error
	// Detail Retrieve the current execution results of a workflow task based on the workflow execution ID.
	Detail(ctx context.Context, workflowRunId string) (*workflowv1.Detail, error)
	// Stop the execution of a workflow task based on the workflow execution ID.
	Stop(ctx context.Context, taskId string) (*workflowv1.StopResponse, error)
	// Logs Returns workflow logs, with the first page returning the latest {limit} messages, i.e., in reverse order.
	Logs(ctx context.Context, filter workflowv1.LogsRequest) (*workflowv1.LogsResponse, error)
}

type workflow struct {
	rest   *resty.Client
	appKey string
}

func (w *workflow) Run(ctx context.Context, request workflowv1.RunRequest) (resp *workflowv1.RunBlockingResponse, err error) {
	_, err = w.rest.R().
		WithContext(ctx).
		SetHeader("x-api-key", w.appKey).
		SetBody(request).
		SetResult(&resp).
		Post("/workflows/run")
	if err != nil {
		return nil, errors.ErrUnsupported
	}
	if resp == nil {
		return nil, errors.ErrUnsupported
	}
	return resp, nil
}

func (w *workflow) RunStream(ctx context.Context, request workflowv1.RunRequest, fn func(v any)) error {
	reqBytes, err := json.Marshal(request)
	if err != nil {
		return err
	}
	es := resty.NewEventSource().
		SetHeader("x-api-key", w.appKey).
		SetMethod(http.MethodPost).
		SetBody(bytes.NewBuffer(reqBytes)).
		OnMessage(func(a any) {
			fn(a)
		}, nil)
	err = es.Get()
	if err != nil {
		return err
	}
	return nil
}

func (w *workflow) Detail(ctx context.Context, workflowRunId string) (*workflowv1.Detail, error) {
	//TODO implement me
	panic("implement me")
}

func (w *workflow) Logs(ctx context.Context, filter workflowv1.LogsRequest) (*workflowv1.LogsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *workflow) Stop(ctx context.Context, taskId string) (*workflowv1.StopResponse, error) {
	//TODO implement me
	panic("implement me")
}

func NewWorkflow(rest *resty.Client, appKey string) Workflow {
	return &workflow{rest: rest, appKey: appKey}
}
