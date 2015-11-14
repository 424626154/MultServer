package LoginServer

import (
	"fmt"
	"lanstonetech.com/network"
	"lanstonetech.com/packet"
	"lanstonetech.com/packet/ID"
)

func ShakeHand(obj *network.SocketBase, msg network.Message) int {

	fmt.Printf("ShakeHand......\n")
	message := new(packet.C2M_Req_ShakeHand)
	err := message.UnPack(msg)
	if err != nil {
		SendShakeHand(obj, ID.REQ_INVALID)
		fmt.Printf("[ShakeHand] UnPack failed! err=%v\n", err)
		return ID.MESSAGE_OK
	}

	fmt.Printf("[ShakeHand] %s\n", message.Greeting)

	SendShakeHand(obj, ID.REQ_SHAKEHAND_OK)
	return ID.MESSAGE_OK
}

func SendShakeHand(obj *network.SocketBase, result uint16) {
	msg := new(packet.M2C_Resp_ShakeHand)
	msg.Result = result

	if message, err := msg.Pack(); err != nil {
		obj.SendMsg(message)
	}
}
