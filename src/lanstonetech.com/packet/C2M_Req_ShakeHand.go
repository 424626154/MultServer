package packet

import (
	"fmt"
	"lanstonetech.com/common"
	"lanstonetech.com/network"
)

type C2M_Req_ShakeHand struct {
	network.CommonPackage

	Greeting string
}

func (this *C2M_Req_ShakeHand) UnPack(msg network.Message) error {
	pos := msg.ParseCommonPackage()

	GreetingLen := common.ReadUint16(msg.Data[pos : pos+2])
	if GreetingLen == 0 || GreetingLen > common.MAX_GREETING_LEN {
		return fmt.Errorf("UnPack failed! GreetingLen==[0, %v]", common.MAX_GREETING_LEN)
	}
	pos += 2
	this.Greeting = common.ReadString(msg.Data[pos:GreetingLen])
	pos += common.MAX_GREETING_LEN

	// this.Greeting = "Hello LoginServer!" //test
	return nil
}
