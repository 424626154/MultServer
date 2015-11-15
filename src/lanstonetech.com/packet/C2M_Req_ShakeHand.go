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
	CommonPackage, pos, err := msg.ParseCommonPackage()
	if err != nil {
		return fmt.Errorf("ParseCommonPackage failed! Err=[%s]", err)
	}
	this.CommonPackage = CommonPackage

	GreetingLen := common.ReadUint16(msg.Data[pos : pos+2])
	if GreetingLen == 0 || GreetingLen > common.MAX_GREETING_LEN {
		return fmt.Errorf("UnPack failed! GreetingLen==[0, %v]", common.MAX_GREETING_LEN)
	}
	pos += 2
	this.Greeting = common.ReadString(msg.Data[pos : pos+int(GreetingLen)])
	pos += common.MAX_GREETING_LEN

	return nil
}
