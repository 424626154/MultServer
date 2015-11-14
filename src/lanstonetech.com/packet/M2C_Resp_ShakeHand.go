package packet

import (
	"lanstonetech.com/common"
	"lanstonetech.com/network"
	"lanstonetech.com/packet/ID"
)

type M2C_Resp_ShakeHand struct {
	Result uint16
}

func (this *M2C_Resp_ShakeHand) MaxLen() uint32 {
	size := 2

	return uint32(size)
}

func (this *M2C_Resp_ShakeHand) Pack() (*network.Message, error) {
	size := this.MaxLen()
	msg := network.NewMessage(uint32(ID.M2C_Resp_ShakeHand), size)
	msg.Data = make([]byte, size)

	pos := 0
	common.WriteUint16(msg.Data[pos:pos+2], this.Result)
	pos += 2

	return msg, nil
}
