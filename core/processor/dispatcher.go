package processor

import (
	"bytes"
	"camellia/core/datapack"
	"context"
	"io/ioutil"
	"net/http"
	"time"
)

type HttpDispatchProcess struct {
	remote  string
	timeout time.Duration

	cli *http.Client
}

func NewHttpDispatcher(remote string, timeout time.Duration) *HttpDispatchProcess {
	d := HttpDispatchProcess{
		remote: remote,
	}

	cli := http.DefaultClient
	cli.Timeout = timeout

	d.cli = cli
	return &d
}

func (d *HttpDispatchProcess) Process(ctx context.Context, msg datapack.Message) *ProcessResp {
	resp, err := d.cli.Post(d.remote, "application/octet-stream", bytes.NewReader(msg.GetPayload()))

	if err != nil {
		return &ProcessResp{Status: 500, err: err}
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &ProcessResp{Status: 500, err: err}
	}
	r := ProcessResp{}
	r.Status = resp.StatusCode
	r.Content = body
	return &r
}
