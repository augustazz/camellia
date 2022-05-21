package processor

import (
	"camellia/core/datapack"
	pb "camellia/pb_generate"
	"context"
)

type Processor interface {
	Process(ctx context.Context, msg datapack.Message) *ProcessResp
}

var defaultProcessor *HttpDispatchProcess

//
var processors = map[pb.MsgType]Processor{
	//default processor
	pb.MsgType_MsgTypeMin: defaultProcessor,
}

type ProcessResp struct {
	Status  int
	err     error
	Content interface{}
}

func GetProcessor(header *pb.Header) Processor {
	if header == nil {
		return nil
	}

	//query by msg type
	msgType := header.MsgType
	if p, ok := processors[msgType]; ok {
		return p
	}

	//query by dest node; pass

	//query default
	if p, ok := processors[0]; ok {
		//index 0 is default processor
		return p
	}

	return nil
}
