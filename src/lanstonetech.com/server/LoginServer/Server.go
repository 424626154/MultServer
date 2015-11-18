package LoginServer

import (
	"lanstonetech.com/common/logger"
	"lanstonetech.com/network"
	"lanstonetech.com/packet/ID"
	"lanstonetech.com/system/config"
	"net"
	"time"
)

var ServerType = 1

func HandlerPackageFunc() {
	network.AddHandler(ID.C2M_Req_ShakeHand, ShakeHand)
}

func InitConf() {
	config.LoadServerInfo(ServerType)
	logger.Infof("[LoginServer] =>> ip=%v port=%v group=%v", config.SERVER_IP, config.SERVER_PORT, config.SERVER_GROUP)
}

func InitLog() {
	logger.SetConsole(true)
	logger.SetConsolePrefix("LoginServer")
	// logger.Initialize("./log", "LoginServer_0")
	logger.SetLevel(logger.LEVEL(2))
}

func Run() {
	defer logger.CatchException()

	InitLog()
	InitConf()
	logger.Infof("server start...")

	HandlerPackageFunc()

	tcpaddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8888")
	if err != nil {
		logger.Errorf("net.ResolveTCPAddr failed! err=%v", err)
		return
	}

	listener, err := net.ListenTCP("tcp4", tcpaddr)
	if err != nil {
		logger.Errorf("net.Listen failed! err=%v", err)
		return
	}

	for {
		conn, err := listener.AcceptTCP()
		logger.Infof("=>%v connecting...", conn.RemoteAddr().String())
		if err != nil {
			logger.Errorf("listener.Accept failed! err=%v", err)
			continue
		}

		go ProcessConnection(conn)
	}

}

func ProcessConnection(conn *net.TCPConn) {
	defer logger.CatchException()

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
			logger.Errorf("SocketBase.RecvMsgs failed! err=%v", err)
			return
		}

		if len(msgs) == 0 {
			continue
		}

		for _, msg := range msgs {
			ret := network.Dispatcher(SocketBase, msg)
			if ret == ID.MESSAGE_OK {
				return
			}
		}
	}
}
