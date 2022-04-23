package channel

import (
	"camellia/core/datapack"
	"camellia/core/enums"
	"camellia/core/event"
	"camellia/core/util"
	"camellia/logger"
	pb "camellia/pb_generate"
	"github.com/golang/protobuf/proto"
	"time"
)


//type handlerFunc func (ctx *ConnContext, msg datapack.Message)
//wrapped by HandlerContext

//HeadDataHandlerFunc head
func HeadDataHandlerFunc(ctx *ConnContext, msg datapack.Message) {
	logger.Debug("head in")

	fail := checkMsg(msg)
	if fail != nil {
		logger.Warning("decode msg is invalid,reason: ", fail.Error())
		ctx.Abort = true
		return
	}

	ctx.LastReadTime = time.Now()
}

//TailDataHandlerFunc tail
func TailDataHandlerFunc(ctx *ConnContext, msg datapack.Message) {
	logger.Info("tail out")
}

//AuthHandlerFunc server verify auth request
func AuthHandlerFunc(ctx *ConnContext, msg datapack.Message) {
	header := msg.GetHeader()
	user := header.UserInfo

	if ctx.Key == "" {
		k := user.Uid
		if k == "" {
			k = user.Did
		}
		ctx.Key = k
	}

	if header.MsgType == pb.MsgType_MsgTypeAuthVerifyReq {
		if ctx.State != enums.ConnStateInAuth {
			event.PostEvent(event.EventTypeConnStatusChanged, ctx.State)
			ctx.State = enums.ConnStateInAuth
		}

		var payload pb.SimplePayload
		err := proto.Unmarshal(msg.GetPayload(), &payload)
		if err != nil {
			logger.Info("err", err)
			return
		}

		succ := verifySig(user, ctx.RandomStr, payload.Content)
		code := pb.AuthCode_AuthFailure
		if succ {
			ctx.State = enums.ConnStateReady
			code = pb.AuthCode_AuthSuccess
		}

		resp := datapack.NewPbMessage()
		resp.HeaderPb.MsgType = pb.MsgType_MsgTypeAuthVerifyResp
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
		logger.Warning("msg not impl processor")
	}
}


func verifySig(user *pb.UserInfo, randomStr string, sig []byte) bool {
	key:= util.GetPubRsaKey()
	if key == nil {
		logger.Warning("get key fail")
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

func checkMsg(msg datapack.Message) error {
	if msg == nil {
		return enums.MsgValidateErrEmpty
	}
	h := msg.GetHeader()
	if h == nil {
		return enums.MsgValidateErrHeaderEmpty
	}
	if h.MsgType == 0 {
		return enums.MsgValidateErrMsgTypeEmpty
	}
	if h.GetUserInfo() == nil {
		return enums.MsgValidateErrUserInfoEmpty
	}

	return nil
}

