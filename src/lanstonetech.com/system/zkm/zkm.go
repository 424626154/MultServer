package zkm

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"lanstonetech.com/system/config"
	"sync"
	"time"
)

var ZKConn zkconn

type zkconn struct {
	sync.RWMutex
	Conn    *zk.Conn
	Root    string
	objects map[string]ZKObject
}

type EventArgs struct {
	Event zk.Event
}

type ZKObject interface {
	Nodes() ([]string, error)
	Request() error
	ProcessException(event *EventArgs) error
	ProcessEvent(event *EventArgs) error
}

func Init() {
	servers := make([]string, 0)

	count := config.ServerConfig.MustInt("ZK", "Count")
	for i := 0; i < count; i++ {
		server := config.ServerConfig.MustValue("ZK", fmt.Sprintf("IP_%d", i))
		port := config.ServerConfig.MustInt("ZK", fmt.Sprintf("PORT_%d", i))
		servers = append(servers, fmt.Sprintf("%s:%d", server, port))
	}

	conn, _, err := Connect(servers)
	if err != nil {
		panic(err)
	}

	root := config.ServerConfig.MustValue("ZK", "Root")

	fmt.Printf("root=%v\n", root)
	ZKConn.objects = make(map[string]ZKObject, 0)
	ZKConn.Root = root
	ZKConn.Conn = conn
}

func Connect(servers []string) (*zk.Conn, <-chan zk.Event, error) {
	tick := config.ServerConfig.MustInt("ZK", "Tick")
	return zk.Connect(servers, time.Duration(tick)*time.Second)
}

func Root() string {
	ZKConn.Lock()
	defer ZKConn.Unlock()

	if ZKConn.Conn == nil {
		panic("[ERR] zk.Conn == nil")
	}

	return ZKConn.Root
}

func AddObserver(node string, obj ZKObject) {
	ZKConn.Lock()
	defer ZKConn.Unlock()

	if len(node) == 0 {
		return
	}

	ZKConn.objects[node] = obj
}

func AddObservers(obj ZKObject) {
	nodes, err := obj.Nodes()
	if err != nil {
		panic(err)
	}

	for _, node := range nodes {
		AddObserver(node, obj)
	}
}

func Start() {
	processRequest()
}

func processRequest() {
	for _, obj := range ZKConn.objects {
		obj.Request()
	}
}

func Create(node string, data string, temp bool) (string, error) {
	ZKConn.Lock()
	defer ZKConn.Unlock()

	if temp {
		return ZKConn.Conn.Create(node, []byte(data), zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	} else {
		return ZKConn.Conn.Create(node, []byte(data), 0, zk.WorldACL(zk.PermAll))
	}
}

func CreateIfNotExists(node string, data string, temp bool) (bool, *zk.Stat, error) {

	if ZKConn.Conn == nil {
		return false, nil, fmt.Errorf("zk.CreateIfNotExists failed! conn == nil")
	}

	exist, stat, err := Exists(node)
	if err != nil {
		return false, nil, err
	}

	if !exist {
		Create(node, data, temp)
	}

	return exist, stat, nil
}

func Exists(node string) (bool, *zk.Stat, error) {
	ZKConn.Lock()
	defer ZKConn.Unlock()

	if ZKConn.Conn == nil {
		return false, nil, fmt.Errorf("zk.Exists failed! conn == nil")
	}

	exist, stat, err := ZKConn.Conn.Exists(node)
	if err != nil {
		return exist, stat, err
	}

	return exist, stat, nil
}

func Get(node string) (string, *zk.Stat, error) {
	ZKConn.Lock()
	defer ZKConn.Unlock()

	if ZKConn.Conn == nil {
		return "", nil, fmt.Errorf("zk.Get failed! conn == nil")
	}

	data, stat, err := ZKConn.Conn.Get(node)
	if err != nil {
		return "", nil, err
	}

	return string(data), stat, nil
}

func Set(node string, data []byte, version int32) (*zk.Stat, error) {
	ZKConn.Lock()
	defer ZKConn.Unlock()

	if ZKConn.Conn == nil {
		return nil, fmt.Errorf("zk.Set failed! conn == nil")
	}

	stat, err := ZKConn.Conn.Set(node, data, version)
	if err != nil {
		return stat, err
	}

	return stat, nil
}

func Delete(node string, version int32) error {
	ZKConn.Lock()
	defer ZKConn.Unlock()

	if ZKConn.Conn == nil {
		return fmt.Errorf("zk.Delete failed! conn == nil")
	}

	return ZKConn.Conn.Delete(node, version)
}

func Children(node string) ([]string, *zk.Stat, error) {
	ZKConn.Lock()
	defer ZKConn.Unlock()

	if ZKConn.Conn == nil {
		return nil, nil, fmt.Errorf("zk.Children failed! conn == nil")
	}

	childrens, stat, err := ZKConn.Conn.Children(node)
	if err != nil {
		return nil, nil, err
	}

	return childrens, stat, nil
}

func GetW(node string) (string, *zk.Stat, error) {
	ZKConn.Lock()
	defer ZKConn.Unlock()

	if ZKConn.Conn == nil {
		return "", nil, fmt.Errorf("zk.GetW failed! conn == nil")
	}

	data, stat, ch, err := ZKConn.Conn.GetW(node)
	if err != nil {
		return "", nil, err
	}

	go func() {
		event := <-ch
		var eventArgs EventArgs
		eventArgs.Event = event

		handlerEvent(&eventArgs)
	}()

	return string(data), stat, nil
}

func ExistsW(node string) (bool, *zk.Stat, error) {
	ZKConn.Lock()
	defer ZKConn.Unlock()

	if ZKConn.Conn == nil {
		return false, nil, fmt.Errorf("zk.ExistsW failed! conn == nil")
	}

	exist, stat, ch, err := ZKConn.Conn.ExistsW(node)
	if err != nil {
		return false, nil, err
	}

	go func() {
		event := <-ch
		var eventArgs EventArgs
		eventArgs.Event = event

		handlerEvent(&eventArgs)
	}()

	return exist, stat, nil
}

func ChildrenW(node string) ([]string, *zk.Stat, error) {
	ZKConn.Lock()
	defer ZKConn.Unlock()

	if ZKConn.Conn == nil {
		return nil, nil, fmt.Errorf("zk.GetW failed! conn == nil")
	}

	data, stat, ch, err := ZKConn.Conn.ChildrenW(node)
	if err != nil {
		return nil, nil, err
	}

	go func() {
		event := <-ch
		var eventArgs EventArgs
		eventArgs.Event = event

		handlerEvent(&eventArgs)
	}()

	return data, stat, nil
}

func handlerEvent(event *EventArgs) {
	if event.Event.State == zk.StateUnknown ||
		event.Event.State == zk.StateDisconnected ||
		event.Event.State == zk.StateExpired ||
		event.Event.State == zk.StateAuthFailed {
		processException(event)
	} else {
		processEvent(event)
	}
}

func processException(event *EventArgs) {
	for node, obj := range ZKConn.objects {
		if node == event.Event.Path {
			go obj.ProcessException(event)
		}
	}
}

func processEvent(event *EventArgs) {
	obj, ok := ZKConn.objects[event.Event.Path]
	if ok {
		obj.ProcessEvent(event)
	}
}
