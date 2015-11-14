package LoginServer

import (
	"fmt"
	"lanstonetech.com/network"
	"lanstonetech.com/packet/ID"
	"net"
	"time"
)

func HandlerPackageFunc() {
	network.AddHandler(ID.C2M_Req_ShakeHand, ShakeHand)
}

func Run() {
	fmt.Printf("server start...\n")

	HandlerPackageFunc()

	tcpaddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Printf("net.ResolveTCPAddr failed! err=%v\n", err)
		return
	}

	listener, err := net.ListenTCP("tcp4", tcpaddr)
	if err != nil {
		fmt.Printf("net.Listen failed! err=%v\n", err)
		return
	}

	for {
		conn, err := listener.AcceptTCP()
		fmt.Printf("=>%v connecting...\n", conn.RemoteAddr().String())
		if err != nil {
			fmt.Printf("listener.Accept failed! err=%v\n", err)
			continue
		}

		go ProcessConnection(conn)
	}

}

func ProcessConnection(conn *net.TCPConn) {
	SocketBase := network.NewSocketBase(conn)
	defer SocketBase.Close()

	SocketBase.SetNoDelay(true)                                        //无延迟
	SocketBase.SetKeepAlive(true)                                      //保持激活
	SocketBase.SetReadBuffer(64 * 1024)                                //设置读缓冲区大小
	SocketBase.SetWriteBuffer(64 * 1024)                               //设置写缓冲区大小
	SocketBase.SetReadDeadline(time.Now().Add(30000000 * time.Second)) //设置读超时

	for {
		msgs, err := SocketBase.RecvMsgs()
		if err != nil {
			fmt.Printf("SocketBase.RecvMsgs failed! err=%v\n", err)
			return
		}

		for _, msg := range msgs {
			ret := network.Dispatcher(SocketBase, msg)
			if ret == ID.MESSAGE_OK {
				return
			}
		}
	}
}
