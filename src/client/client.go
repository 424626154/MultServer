package main

import (
	"fmt"
	"lanstonetech.com/common"
	"lanstonetech.com/common/logger"
	"lanstonetech.com/network"
	"net"
)

var (
	account   = "dcea85c06887c52f20994c393c5bbcdc"
	signature = "dcea85c06887c52f20994c393c5bbcdc"
	token     = "dcea85c06887c52f20994c393c5bbcdc"
	greeting  = "Hello LoginServer"
)

//====================================================================
//------------------Send greeting to server --------------------------
//====================================================================

func main() {
	//====================================================================
	//---------------------------- Data ----------------------------------
	//====================================================================
	size := uint32(common.MAX_HEADER_LEN + common.COMMON_PACKAGE_LEN + 2 + len(greeting))
	header := new(network.Header)
	header.PacketID = 100
	header.PackageLen = size - 8

	//====================================================================
	//---------------------------- Pack Message---------------------------
	//====================================================================
	message := network.NewMessage(100, size-common.MAX_HEADER_LEN)
	message.Data = make([]byte, size)

	//CommonPackage
	message.Account = account
	message.Signature = signature
	message.Token = token

	//Pack Message
	pos := message.PackHeader()
	pos = message.PackCommonPackage(pos)

	//Pach body
	leng := len(greeting)
	common.WriteUint16(message.Data[pos:pos+2], uint16(leng))
	pos += 2
	common.WriteString(message.Data[pos:pos+leng], greeting)
	pos += common.MAX_GREETING_LEN

	//====================================================================
	//---------------------------- New Session ---------------------------
	//====================================================================
	tcpaddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:4000")
	if err != nil {
		logger.Errorf("[client] net.ResolveTCPAddr failed! err=%v", err)
	}

	conn, err := net.DialTCP("tcp", nil, tcpaddr)
	if err != nil {
		logger.Errorf("[client] net.DiaTCP failed! err=%v", err)
	}
	defer func() {
		conn.Close()
		logger.Errorf("[client] ok...!")
	}()

	//====================================================================
	//---------------------------- Send ----------------------------------
	//====================================================================
	TCPConn := network.NewSocketBase(conn)
	leng, err = TCPConn.Write(message.Data)
	if err != nil {
		logger.Errorf("[client] conn.Write failed! err=%v", err)
	}

	//====================================================================
	//---------------------------- Recv ----------------------------------
	//====================================================================
	//Parse Message
	msgs, err := TCPConn.RecvMsgs()
	if err != nil {
		logger.Errorf("err=%v", err)
		return
	}

	//handler messages
	for _, msg := range msgs {
		result := common.ReadUint16(msg.Data[0:2])
		logger.Errorf("result=%v\n", result)
	}

	return
}
