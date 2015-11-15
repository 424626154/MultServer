package network

import (
	"fmt"
	"net"
	"sync"
)

type SocketBase struct {
	sync.RWMutex
	*net.TCPConn
	// Param SocketBase
	senddata  chan Message
	send_flag chan int
}

func NewSocketBase(conn *net.TCPConn) *SocketBase {
	socket := &SocketBase{
		TCPConn: conn,
	}

	return socket
}

func (this *SocketBase) IsPacket(id int32) bool {
	return true
}

func (this *SocketBase) RecvMsgs() ([]Message, error) {
	messages := make([]Message, 0)

	// for {
	var message Message
	header_buf := make([]byte, MAX_HEADER_LEN)

	//读取头部
	length, err := this.TCPConn.Read(header_buf)
	if err != nil {
		return nil, err
	}

	message.ParseHeader(header_buf)

	//读取包正文
	body_buf := make([]byte, message.PackageLen)
	length, err = this.TCPConn.Read(body_buf)
	if err != nil {
		return nil, err
	}

	if length > MAX_PACKAGE_LEN {
		return nil, fmt.Errorf("[ERR] Packet size is too big")
	}

	message.Data = make([]byte, length)
	copy(message.Data, body_buf)

	messages = append(messages, message)
	// }

	return messages, nil
}

func (this *SocketBase) SendMsg(msg *Message) error {
	this.RWMutex.Lock()
	defer func() {
		this.RWMutex.Unlock()
	}()

	send_data := make([]byte, 4+4+msg.PackageLen)
	header_len := msg.WriteHeader(send_data)

	copy(send_data[header_len:], msg.Data)

	_, err := this.TCPConn.Write(send_data)
	if err != nil {
		return err
	}

	return nil
}

func (this *SocketBase) Close() {
	this.RWMutex.Lock()
	defer func() {
		this.RWMutex.Unlock()
		this.TCPConn.Close()
	}()
}
