package channel

import (
	"camellia/core/datapack"
	"camellia/core/enums"
	"camellia/core/event"
	"camellia/core/util"
	pb "camellia/pb_generate"
	"fmt"
	"github.com/golang/protobuf/proto"
)


//type handlerFunc func (ctx *ConnContext, msg datapack.Message)
//wrapped by HandlerContext

//HeadDataHandlerFunc head
func HeadDataHandlerFunc(ctx *ConnContext, msg datapack.Message) {
	fmt.Println("head in")
}

//TailDataHandlerFunc tail
func TailDataHandlerFunc(ctx *ConnContext, msg datapack.Message) {
	fmt.Println("tail out")
}

//AuthHandlerFunc server verify auth request
func AuthHandlerFunc(ctx *ConnContext, msg datapack.Message) {
	if msg.GetHeader().MsgType == pb.MsgType_MsgTypeAuthVerifyReq {
		if ctx.State != enums.ConnStateInAuth {
			event.PostEvent(event.EventTypeConnStatusChanged, ctx.State)
			ctx.State = enums.ConnStateInAuth
		}

		var payload pb.SimplePayload
		err := proto.Unmarshal(msg.GetPayload(), &payload)
		if err != nil {
			fmt.Println("err", err)
			return
		}

		user := msg.GetHeader().UserInfo
		if user == nil {
			return
		}
		succ := verifySig(msg.GetHeader().UserInfo, ctx.RandomStr, payload.Content)

		code := pb.AuthCode_AuthFailure
		if succ {
			ctx.State = enums.ConnStateReady
			code = pb.AuthCode_AuthSuccess
		}

		resp := datapack.NewPbMessage()
		resp.Header.MsgType = pb.MsgType_MsgTypeAuthVerifyResp
		resp.PayloadPb = &pb.AuthResp{
			Code: code,
		}
		ctx.WriteChan<- (&datapack.TcpPackage{}).Pack(resp)
		return
	}

	//没验证通过时，数据丢弃
	if ctx.State != enums.ConnStateReady {
		ctx.Abort = true
	}
}

var msgProcessors map[pb.MsgType]func(datapack.Message)

func DispatchHandlerFunc(ctx *ConnContext, msg datapack.Message) {
	msgType := msg.GetHeader().MsgType
	processor, ok := msgProcessors[msgType]
	if !ok {
		//if has default processor
		processor, ok = msgProcessors[0]
	}

	if ok {
		processor(msg)
	} else {
		fmt.Println("msg not impl processor")
	}
}


func verifySig(user *pb.UserInfo, randomStr string, sig []byte) bool {
	key:= util.GetPubRsaKey()
	if key == nil {
		fmt.Println("get key fail")
		return false
	}
	uid := user.Uid
	did := user.Did

	content := make([]byte, 0, len(uid) + len(did) + len(randomStr))
	content = append(content, []byte(uid)...)
	content = append(content, []byte(did)...)
	content = append(content, []byte(randomStr)...)

	return util.RsaVerySignWithSha256(content, sig, key)
}

