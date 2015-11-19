package network

import (
	"net"
	"sync"
)

type SocketBase struct {
	sync.RWMutex
	*net.TCPConn
	// Param SocketBase
	senddata  chan Message
	send_flag chan int

	buffer []byte //buffer for receive
	pos    int    //index for write
}

func NewSocketBase(conn *net.TCPConn) *SocketBase {
	socket := &SocketBase{
		TCPConn: conn,
		buffer:  make([]byte, 0), //the max size of buffer
		pos:     0,
	}

	return socket
}

func (this *SocketBase) IsPacket(id int32) bool {
	return true
}

func (this *SocketBase) RecvMsgs() ([]Message, error) {
	messages := make([]Message, 0)

	//读数据
	recv_data := make([]byte, MAX_PACKAGE_LEN)
	n, err := this.TCPConn.Read(recv_data)
	if err != nil {
		return nil, err
	}

	if n == 0 {
		return nil, nil
	}

	//update buffer
	temp := this.buffer[:this.pos]
	this.buffer = append(temp, recv_data[:n]...)
	this.pos += n

	for {
		if this.pos < MAX_HEADER_LEN {
			break
		}

		var message Message

		//parse header
		message.ParseHeader(this.buffer[0:MAX_HEADER_LEN])

		if this.pos < MAX_HEADER_LEN+int(message.PackageLen) {
			break
		}

		//read body
		pos := MAX_HEADER_LEN + int(message.PackageLen)
		message.Data = this.buffer[MAX_HEADER_LEN:pos]

		//handle remain buffer
		temp := this.buffer[pos:this.pos]
		this.buffer = temp
		this.pos -= pos

		messages = append(messages, message)
	}

	return messages, nil
}

//channel
func (this *SocketBase) SendMsg(msg *Message) error {
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
	this.TCPConn.Close()
	this.buffer = nil
	this.pos = 0
}
