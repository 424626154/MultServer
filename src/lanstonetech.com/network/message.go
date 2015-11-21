package network

import (
	"fmt"
	"lanstonetech.com/common"
	// "net"
)

const (
	MAX_PACKAGE_LEN    = 1024
	MAX_HEADER_LEN     = 4 + 4
	MAX_ACCOUNT_LEN    = 48
	MAX_SIGNATURE_LEN  = 48
	MAX_TOKEN_LEN      = 48
	COMMON_PACKAGE_LEN = 2 + MAX_ACCOUNT_LEN + 2 + MAX_SIGNATURE_LEN + 2 + MAX_TOKEN_LEN
)

type Message struct {
	Header
	CommonPackage
	Data []byte
}

type Header struct {
	PacketID   uint32
	PackageLen uint32
}

type CommonPackage struct {
	Account   string
	Signature string
	Token     string
}

//=============================================================================
//------------------------- Server Action -------------------------------------
//=============================================================================
func NewMessage(packet_id uint32, size uint32) *Message {
	msg := new(Message)
	msg.PacketID = packet_id
	msg.PackageLen = size

	return msg
}

//解析协议头
func (this *Message) ParseHeader(buf []byte) {
	pos := 0
	this.PacketID = common.ReadUint32(buf[pos : pos+4])
	pos += 4
	this.PackageLen = common.ReadUint32(buf[pos : pos+4])
	pos += 4

	this.Data = make([]byte, this.PackageLen)
}

//解析公共包头
func (this *Message) ParseCommonPackage() (CommonPackage, int, error) {

	pos := uint16(0)
	AccountLen := common.ReadUint16(this.Data[pos : pos+2])
	if AccountLen == 0 || AccountLen > MAX_ACCOUNT_LEN {
		return CommonPackage{}, 0, fmt.Errorf("UnPack failed! AccountLen==[0, %v]", MAX_ACCOUNT_LEN)
	}
	pos += 2
	account := common.ReadString(this.Data[pos : pos+AccountLen])
	pos += MAX_ACCOUNT_LEN

	SignatureLen := common.ReadUint16(this.Data[pos : pos+2])
	if SignatureLen == 0 || SignatureLen > MAX_SIGNATURE_LEN {
		return CommonPackage{}, 0, fmt.Errorf("UnPack failed! SignatureLen==[0, %v]", MAX_SIGNATURE_LEN)
	}
	pos += 2
	signature := common.ReadString(this.Data[pos : pos+SignatureLen])
	pos += MAX_SIGNATURE_LEN

	TokenLen := common.ReadUint16(this.Data[pos : pos+2])
	if TokenLen == 0 || TokenLen > MAX_TOKEN_LEN {
		return CommonPackage{}, 0, fmt.Errorf("UnPack failed! TokenLen==[0, %v]", MAX_TOKEN_LEN)
	}
	pos += 2
	token := common.ReadString(this.Data[pos : pos+TokenLen])
	pos += MAX_TOKEN_LEN

	var common_package CommonPackage
	common_package.Account = account
	common_package.Signature = signature
	common_package.Token = token

	return common_package, int(pos), nil
}

func (this *Message) PackMessage() []byte {
	buf := make([]byte, MAX_HEADER_LEN+this.PackageLen)

	pos := 0
	common.WriteUint32(buf[pos:pos+4], this.PacketID)
	pos += 4
	common.WriteUint32(buf[pos:pos+4], this.PackageLen)
	pos += 4

	copy(buf[pos:], this.Data[0:this.PackageLen])

	return buf
}

//=============================================================================
//------------------------- Client Action -------------------------------------
//=============================================================================
func (this *Message) PackHeader() int {
	pos := 0
	common.WriteUint32(this.Data[pos:pos+4], this.PacketID)
	pos += 4
	common.WriteUint32(this.Data[pos:pos+4], this.PackageLen)
	pos += 4

	return pos
}

func (this *Message) PackCommonPackage(index int) int {
	pos := index

	leng := len(this.Account)
	common.WriteUint16(this.Data[pos:pos+2], uint16(leng))
	pos += 2
	common.WriteString(this.Data[pos:pos+leng], this.Account)
	pos += 48

	leng = len(this.Signature)
	common.WriteUint16(this.Data[pos:pos+2], uint16(leng))
	pos += 2
	common.WriteString(this.Data[pos:pos+leng], this.Signature)
	pos += 48

	leng = len(this.Token)
	common.WriteUint16(this.Data[pos:pos+2], uint16(leng))
	pos += 2
	common.WriteString(this.Data[pos:pos+leng], this.Token)
	pos += 48

	return pos
}
