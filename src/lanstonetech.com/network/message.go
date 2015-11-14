package network

import (
	// "fmt"
	"lanstonetech.com/common"
	// "net"
)

const (
	MAX_PACKAGE_LEN    = 2048
	MAX_HEADER_LEN     = 4 + 4
	MAX_ACCOUNT_LEN    = 48
	MAX_SIGNATURE_LEN  = 48
	MAX_TOKEN_LEN      = 48
	COMMON_PACKAGE_LEN = MAX_ACCOUNT_LEN + MAX_SIGNATURE_LEN + MAX_TOKEN_LEN
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
}

//解析协议头
func (this *Message) WriteHeader(buf []byte) int {
	pos := 0
	common.WriteUint32(buf[pos:pos+4], this.PacketID)
	pos += 4
	common.WriteUint32(buf[pos:pos+4], this.PackageLen)
	pos += 4

	return MAX_HEADER_LEN
}

//解析公共包头
func (this *Message) ParseCommonPackage() int {
	pos := 0
	this.Account = common.ReadString(this.Data[pos : pos+MAX_ACCOUNT_LEN])
	pos += MAX_ACCOUNT_LEN

	this.Signature = common.ReadString(this.Data[pos : pos+MAX_SIGNATURE_LEN])
	pos += MAX_SIGNATURE_LEN

	this.Token = common.ReadString(this.Data[pos : pos+MAX_TOKEN_LEN])
	pos += MAX_TOKEN_LEN

	return pos
}
