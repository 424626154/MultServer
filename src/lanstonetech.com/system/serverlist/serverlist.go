package serverlist

import (
	"encoding/json"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"lanstonetech.com/common/logger"
	"lanstonetech.com/system/config"
	"lanstonetech.com/system/zkm"
	"sync"
	"time"
)

var ServerList serverlist

type serverlist struct {
	sync.RWMutex
	// serverGroup string
	serverTypes []string
	servers     map[uint8][]zkm.ZKServerInfo
}

func init() {
	ServerList.servers = make(map[uint8][]zkm.ZKServerInfo)
}

// func Init(group string) {

// }

func (this *serverlist) AddMonitor(serverGroup string, serverType uint8) {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	this.serverTypes = append(this.serverTypes, zkm.Server.MakeServerListGroupNode(serverGroup, serverType))
}

func (this *serverlist) Nodes() ([]string, error) {
	this.RLock()
	defer this.RUnlock()

	return this.serverTypes, nil
}

func (this *serverlist) Request() error {
	this.RequestServerList()

	return nil
}

func (this *serverlist) RequestServerList() {
	for _, item := range this.serverTypes {
		go this.requestServerList(item)
	}
}

func (this *serverlist) requestServerList(serverGroup string) {
	defer logger.CatchException()

	for {
		sl, _, err := zkm.ChildrenW(serverGroup)
		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}

		this.parseServerList(serverGroup, sl)
	}
}

func (this *serverlist) parseServerList(serverGroup string, sl []string) {
	for _, s := range sl {
		var serverInfo zkm.ZKServerInfo
		err := json.Unmarshal([]byte(s), &serverInfo)
		if err != nil {
			logger.Errorf("parseServerList failed! err=%v", err)
			return
		}

		if serverInfo.Class == config.SERVER_TYPE && serverInfo.IP == config.SERVER_IP && serverInfo.Port == config.SERVER_PORT {
			continue
		}

		this.addServer(&serverInfo)
	}
}

func (this *serverlist) addServer(server *zkm.ZKServerInfo) {
	if this.hasServer(server) {
		return
	}

	sl, ok := this.servers[server.Class]
	if !ok {
		temp := []zkm.ZKServerInfo{*server}
		this.servers[server.Class] = temp
	} else {
		sl = append(sl, *server)
		this.servers[server.Class] = sl
	}
}

func (this *serverlist) hasServer(server *zkm.ZKServerInfo) bool {
	sl, ok := this.servers[server.Class]
	if !ok {
		return false
	}

	for _, s := range sl {
		if s.Group == server.Group && s.Class == server.Class && s.IP == server.IP && s.Port == server.Port {
			return true
		}
	}

	return false
}

func (this *serverlist) ProcessException(event *zkm.EventArgs) error {
	return this.Request()
}

func (this *serverlist) ProcessEvent(event *zkm.EventArgs) error {
	if event.Event.Type == zk.EventNodeCreated {
		//节点被创建
	} else if event.Event.Type == zk.EventNodeDeleted {
		//节点删除
	} else if event.Event.Type == zk.EventNodeDataChanged {
		//节点数据改变
	} else if event.Event.Type == zk.EventNodeChildrenChanged {
		//子节点改变
		this.RequestServerList()
	}

	return nil
}

func (this *serverlist) GetServerList(class uint8) ([]zkm.ZKServerInfo, error) {
	this.Unlock()
	defer this.Lock()

	sl, ok := this.servers[class]
	if !ok {
		return []zkm.ZKServerInfo{}, fmt.Errorf("[serverlist] GetServerList failed! have no servers")
	}

	return sl, nil
}
