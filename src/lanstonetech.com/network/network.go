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
		buffer:  make([]byte, 5*MAX_PACKAGE_LEN), //the max size of buffer
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

	//update buffer
	copy(this.buffer[this.pos:this.pos+n], recv_data[0:n])
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
		message.Data = this.buffer[MAX_HEADER_LEN : MAX_HEADER_LEN+int(message.PackageLen)]

		//handle remain buffer
		this.buffer = this.buffer[MAX_HEADER_LEN+int(message.PackageLen):]
		this.pos -= MAX_HEADER_LEN + int(message.PackageLen)

		messages = append(messages, message)
	}

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
