package dify

import "resty.dev/v3"

import "github.com/openexw/dify-go/workflow"

type client struct {
	rest     *resty.Client
	workflow workflow.Workflow
}

type Client interface {
	Workflow(appKey string) workflow.Workflow
}

func NewClient() Client {
	c := resty.New()
	return &client{
		rest: c,
	}
}

func (cli *client) Workflow(appKey string) workflow.Workflow {
	return workflow.NewWorkflow(cli.rest, appKey)
}
