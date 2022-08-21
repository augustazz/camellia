package processor

import (
	"context"
	"github.com/augustazz/camellia/core/datapack"
	pb "github.com/augustazz/camellia/pb_generate"
)

type Processor interface {
	Process(ctx context.Context, msg datapack.Message) *ProcessResp
}

var defaultProcessor *HttpDispatchProcess

//
var processors = map[pb.MsgType]Processor{
	//default processor
	pb.MsgType_Min: defaultProcessor,
}

type ProcessResp struct {
	Success bool
	Finish  bool
	err     error
	Content interface{}
}
