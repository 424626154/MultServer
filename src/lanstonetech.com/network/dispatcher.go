package network

var handlers map[uint32]func(*SocketBase, *Message) int

func init() {
	handlers = make(map[uint32]func(*SocketBase, *Message) int, 0)
}

func AddHandler(PackageID uint32, handler func(*SocketBase, *Message) int) {
	handlers[PackageID] = handler
}

func Dispatcher(conn *SocketBase, msg *Message) int {
	for k, v := range handlers {
		if k == msg.PacketID {
			return v(conn, msg)
		}
	}

	return 0
}
