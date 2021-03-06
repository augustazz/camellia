package handler

import (
	"camellia/core/channel"
	"camellia/core/datapack"
	"camellia/core/enums"
	"camellia/core/util"
	"camellia/logger"
	pb "camellia/pb_generate"
	"github.com/golang/protobuf/proto"
)


func ClientAuthHandlerFunc(ctx *channel.ConnContext, msg datapack.Message) {
	switch msg.GetHeader().MsgType {
	case pb.MsgType_MsgTypeAuthLaunch:
		ctx.State = enums.ConnStateInAuth

		var payload pb.SimplePayload
		err := proto.Unmarshal(msg.GetPayload(), &payload)
		if err != nil {
			logger.Error(err)
			return
		}
		origin := payload.Content
		resp := datapack.NewPbMessage()
		resp.GetHeader().MsgType = pb.MsgType_MsgTypeAuthVerifyReq
		resp.GetHeader().UserInfo = &pb.UserInfo{
			Uid: "100023",
			Did: "DT39485",
		}

		resp.PayloadPb = &pb.SimplePayload{
			Content: encrypt(resp.GetHeader().UserInfo, origin),
		}

		ctx.WriteChan<- (&datapack.TcpPackage{}).Pack(resp)
	case pb.MsgType_MsgTypeAuthVerifyResp:
		var result pb.AuthResp
		err := proto.Unmarshal(msg.GetPayload(), &result)
		if err != nil {
			logger.Error(err)
			return
		}
		if result.Code == pb.AuthCode_AuthSuccess {
			logger.Info("auth success")
			ctx.State = enums.ConnStateReady
		} else {
			logger.Info("auth fail")
		}
	}
}



func encrypt(user *pb.UserInfo, origin []byte) []byte {
	prvKey := util.GetPrvRsaKey()
	if prvKey == nil {
		logger.Info("prvKey fail")
		return nil
	}

	//uid+did+random str
	s := make([]byte, 0, len(user.Uid) + len(user.Did) + len(origin))
	s = append(s, []byte(user.Uid)...)
	s = append(s, []byte(user.Did)...)
	s = append(s, origin...)
	logger.Debug("content:", string(s))
	return util.RsaSignWithSha256(s, prvKey)
}

