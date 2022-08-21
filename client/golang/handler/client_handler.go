package handler

import (
	"github.com/augustazz/camellia/config"
	"github.com/augustazz/camellia/constants"
	"github.com/augustazz/camellia/core/channel"
	"github.com/augustazz/camellia/core/datapack"
	"github.com/augustazz/camellia/logger"
	pb "github.com/augustazz/camellia/pb_generate"
	"github.com/augustazz/camellia/util"
	"github.com/golang/protobuf/proto"
)

func ClientAuthHandlerFunc(ctx *channel.ConnContext, msg datapack.Message) {
	switch msg.GetHeader().MsgType {
	case pb.MsgType_AuthLaunch:
		ctx.State = constants.ConnStateInAuth

		var payload pb.SimpleMessage
		err := proto.Unmarshal(msg.GetPayload(), &payload)
		if err != nil {
			logger.Error(err)
			return
		}
		origin := payload.Content
		if len(origin) == 0 {
			logger.Error("server resp auth content is empty")
			return
		}
		resp := datapack.NewPbMessageWithEndpoint(pb.Endpoint_Client, pb.Endpoint_ServerConnCenter)
		resp.GetHeader().MsgType = pb.MsgType_AuthVerifyReq
		resp.GetHeader().UserInfo = &pb.UserInfo{
			Uid: "100023",
			Did: "DT39485",
		}

		resp.PayloadPb = &pb.SimpleMessage{
			Content: encrypt(resp.GetHeader().UserInfo, origin),
		}

		ctx.WriteChan <- (&datapack.TcpPackage{}).Pack(resp)
	case pb.MsgType_AuthVerifyResp:
		var result pb.AuthRespMessage
		err := proto.Unmarshal(msg.GetPayload(), &result)
		if err != nil {
			logger.Error(err)
			return
		}
		if result.Code == pb.AuthCode_AuthSuccess {
			logger.Info("auth success")
			ctx.State = constants.ConnStateReady
		} else {
			logger.Info("auth fail")
		}
	}
}

func encrypt(user *pb.UserInfo, origin []byte) []byte {
	prvKey := util.GetPrvRsaKey(config.GetConnConfig().AuthFilePath)
	if prvKey == nil {
		logger.Info("prvKey fail")
		return nil
	}

	//uid+did+random str
	s := make([]byte, 0, len(user.Uid)+len(user.Did)+len(origin))
	s = append(s, []byte(user.Uid)...)
	s = append(s, []byte(user.Did)...)
	s = append(s, origin...)
	logger.Debug("content:", string(s))
	return util.RsaSignWithSha256(s, prvKey)
}
