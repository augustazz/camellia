package processor

import (
	"bytes"
	"context"
	"github.com/augustazz/camellia/constants"
	"github.com/augustazz/camellia/core/datapack"
	"github.com/augustazz/camellia/core/discovery"
	"github.com/augustazz/camellia/core/transport"
	pb "github.com/augustazz/camellia/pb_generate"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	clientEndpointProcessor = new(DownstreamProcessor)
	dispatchProcess         = new(HttpDispatchProcess)
)

func PushToProcessor(msg datapack.Message) *ProcessResp {
	var p Processor
	e := msg.GetHeader().Dest
	if e == pb.Endpoint_Client {
		p = clientEndpointProcessor
	} else {
		p = dispatchProcess
	}

	return p.Process(context.Background(), msg)
}

type DownstreamProcessor struct {
}

func (d *DownstreamProcessor) Process(ctx context.Context, msg datapack.Message) *ProcessResp {

	r := ProcessResp{}

	return &r
}

type RemoteDispatchProcess interface {
	Processor

	GetRemoteServices(serviceName string) []transport.RemoteService
	LoadBalanceSort([]transport.RemoteService) []transport.RemoteService
	DoDispatch() *ProcessResp
}

type HttpDispatchProcess struct {
	endpoint pb.Endpoint
	cli      *http.Client
	timeout  time.Duration
}

func NewHttpDispatcher(timeout time.Duration) *HttpDispatchProcess {
	d := HttpDispatchProcess{
		timeout: timeout,
	}

	cli := http.DefaultClient
	cli.Timeout = timeout

	d.cli = cli
	return &d
}

func (d *HttpDispatchProcess) Process(ctx context.Context, msg datapack.Message) *ProcessResp {
	services := discovery.GetDiscovery().GetServices(d.endpoint.String())
	if services == nil || len(services.GetServiceUrl()) == 0 {
		return &ProcessResp{Success: false, err: constants.ServicesNotAvailable}
	}

	resp, err := d.cli.Post(services.GetServiceUrl()[0], "application/octet-stream", bytes.NewReader(msg.GetPayload()))

	if err != nil {
		return &ProcessResp{Success: false, err: err}
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &ProcessResp{Success: false, err: err}
	}
	r := ProcessResp{}
	r.Success = true
	r.Content = body
	return &r
}
