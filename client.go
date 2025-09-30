package dify

import "resty.dev/v3"

type client struct {
	rest     *resty.Client
	workflow Workflow
}

type Client interface {
	Workflow(appKey string) Workflow
}

func NewClient() Client {
	c := resty.New()
	return &client{
		rest: c,
	}
}

func (cli *client) Workflow(appKey string) Workflow {
	return NewWorkflow(cli.rest, appKey)
}
