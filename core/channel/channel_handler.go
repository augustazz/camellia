package channel

import (
	"github.com/augustazz/camellia/config"
	"github.com/augustazz/camellia/constants"
	"github.com/augustazz/camellia/core/datapack"
	"github.com/augustazz/camellia/core/event"
	"github.com/augustazz/camellia/core/processor"
	"github.com/augustazz/camellia/logger"
	pb "github.com/augustazz/camellia/pb_generate"
	"github.com/augustazz/camellia/util"
	"github.com/golang/protobuf/proto"
	"time"
)

//type handlerFunc func (ctx *ConnContext, msg datapack.Message)
//wrapped by HandlerContext

//HeadDataHandlerFunc head
func HeadDataHandlerFunc(ctx *ConnContext, msg datapack.Message) {
	defer func() {
		err := recover()
		if err != nil {
			logger.Error("recover error:%v", err)
			//close client conn? todo
		}
	}()

	logger.Debug("head in")

	fail := checkMsg(msg)
	if fail != nil {
		logger.Warning("decode msg is invalid,reason: ", fail.Error())
		ctx.Abort()
		return
	}

	ctx.LastReadTime = time.Now()
}

//TailDataHandlerFunc tail
func TailDataHandlerFunc(ctx *ConnContext, msg datapack.Message) {
	logger.Info("tail out")
	ctx.AbortReset()
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

	if header.MsgType == pb.MsgType_AuthVerifyReq {
		resp := datapack.NewPbMessageWithEndpoint(pb.Endpoint_ServerConnCenter, pb.Endpoint_Client)
		resp.HeaderPb.MsgType = pb.MsgType_AuthVerifyResp

		if ctx.State != constants.ConnStateInAuth {
			event.PostEvent(event.EventTypeConnStatusChanged, ctx.State)
			ctx.State = constants.ConnStateInAuth
		}

		var payload pb.SimpleMessage
		err := proto.Unmarshal(msg.GetPayload(), &payload)
		if err != nil {
			logger.Errorf("pb unmarshal err", err)
			return
		}
		if len(payload.GetContent()) == 0 {
			logger.Warningf("auth req fail[%s],header:[%s]", "payload is empty", msg.GetHeader().String())
			resp.PayloadPb = &pb.AuthRespMessage{
				Code: pb.AuthCode_AuthFailure,
			}
			ctx.WriteChan <- (&datapack.TcpPackage{}).Pack(resp)
			return
		}

		succ := verifySig(user, ctx.RandomStr, payload.Content)
		code := pb.AuthCode_AuthFailure
		if succ {
			ctx.State = constants.ConnStateReady
			code = pb.AuthCode_AuthSuccess
		}
		resp.PayloadPb = &pb.AuthRespMessage{
			Code: code,
		}
		ctx.WriteChan <- (&datapack.TcpPackage{}).Pack(resp)
		return
	}

	//没验证通过时，数据丢弃
	if ctx.State != constants.ConnStateReady {
		ctx.Abort()
	}
}

func DispatchHandlerFunc(connCtx *ConnContext, msg datapack.Message) {
	h := msg.GetHeader()

	if h.Dest == pb.Endpoint_ServerConnCenter {
		return
	}

	err := processor.PushToProcessor(msg)
	if err != nil {
		logger.Error("msg process err:", err)
	}
}

func AckHandlerFunc(connCtx *ConnContext, msg datapack.Message) {
	h := msg.GetHeader()
	if h.Ack { //ack
		resp := datapack.NewPbMessageWithEndpoint(pb.Endpoint_ServerConnCenter, h.Src)
		resp.GetHeader().MsgType = pb.MsgType_ServerAck
		payload := &pb.AckMessage{
			SourceMsgId:   h.MsgId,
			SourceMsgType: h.MsgType,

			Code:      pb.AckCode_Received,
			Timestamp: uint64(time.Now().UnixNano() / 1e6),
		}
		resp.PayloadPb = payload
		connCtx.WriteChan <- (&datapack.TcpPackage{}).Pack(resp)
	}
}

func verifySig(user *pb.UserInfo, randomStr string, sig []byte) bool {
	key := util.GetPubRsaKey(config.GetSrvConfig().Conn.AuthFilePath)
	if key == nil {
		logger.Warning("get key fail")
		return false
	}
	uid := user.Uid
	did := user.Did

	content := make([]byte, 0, len(uid)+len(did)+len(randomStr))
	content = append(content, []byte(uid)...)
	content = append(content, []byte(did)...)
	content = append(content, []byte(randomStr)...)

	return util.RsaVerySignWithSha256(content, sig, key)
}

func checkMsg(msg datapack.Message) error {
	if msg == nil {
		return constants.MsgValidateErrEmpty
	}
	h := msg.GetHeader()
	if h == nil {
		return constants.MsgValidateErrHeaderEmpty
	}
	if h.MsgType == 0 {
		return constants.MsgValidateErrMsgTypeEmpty
	}
	if h.Src == pb.Endpoint_Client && h.GetUserInfo() == nil {
		return constants.MsgValidateErrUserInfoEmpty
	}

	return nil
}
