package LoginServer

import (
	"lanstonetech.com/common/logger"
	"lanstonetech.com/network"
	"lanstonetech.com/packet"
	"lanstonetech.com/packet/ID"
)

func ShakeHand(obj *network.SocketBase, msg *network.Message) int {

	message := new(packet.C2M_Req_ShakeHand)
	err := message.UnPack(msg)
	if err != nil {
		SendShakeHand(obj, ID.REQ_INVALID)
		logger.Errorf("[ShakeHand] UnPack failed! err=%v\n", err)
		return ID.MESSAGE_OK
	}

	SendShakeHand(obj, ID.REQ_SHAKEHAND_OK)
	logger.Errorf("[ShakeHand] Process successful! Greeting=%s\n", message.Greeting)

	return ID.MESSAGE_OK
}

func SendShakeHand(obj *network.SocketBase, result uint16) {
	msg := new(packet.M2C_Resp_ShakeHand)
	msg.Result = result

	message, err := msg.Pack()
	if err != nil {
		logger.Errorf("msg.Pack failed! err=%v", err)
		return
	}

	if err := obj.SendMsg(message); err != nil {
		logger.Errorf("obj.SendMsg failed! err=%v", err)
	}
}
