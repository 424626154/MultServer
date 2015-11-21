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
	serverTypes []string
	servers     map[uint8][]zkm.ZKServerInfo
}

func init() {
	ServerList.servers = make(map[uint8][]zkm.ZKServerInfo)
}

func (this *serverlist) AddMonitor(serverGroup string, serverType uint8) {
	this.Lock()
	defer this.Unlock()

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
		return
	}
}

func (this *serverlist) parseServerList(serverGroup string, sl []string) {
	cleanup := true
	for _, s := range sl {
		var serverInfo zkm.ZKServerInfo
		if err := json.Unmarshal([]byte(s), &serverInfo); err != nil {
			logger.Errorf("parseServerList failed! err=%v", err)
			return
		}

		if serverInfo.Type == config.SERVER_TYPE && serverInfo.IP == config.SERVER_IP && serverInfo.Port == config.SERVER_PORT {
			continue
		}

		this.addServer(&serverInfo, cleanup)
		if cleanup {
			cleanup = false
		}
	}
}

func (this *serverlist) cleanupServer(serverType uint8) {
	_, ok := this.servers[serverType]
	if ok {
		delete(this.servers, serverType)
	}
}

func (this *serverlist) addServer(server *zkm.ZKServerInfo, cleanup bool) {
	if cleanup {
		this.cleanupServer(server.Type)
	}

	if this.hasServer(server) {
		return
	}

	sl, ok := this.servers[server.Type]
	if !ok {
		temp := []zkm.ZKServerInfo{*server}
		this.servers[server.Type] = temp
	} else {
		sl = append(sl, *server)
		this.servers[server.Type] = sl
	}
}

func (this *serverlist) hasServer(server *zkm.ZKServerInfo) bool {
	sl, ok := this.servers[server.Type]
	if !ok {
		return false
	}

	for _, s := range sl {
		if s.Group == server.Group && s.Type == server.Type && s.IP == server.IP && s.Port == server.Port {
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
		this.requestServerList(event.Event.Path)
	}

	return nil
}

func (this *serverlist) GetServerList(class uint8) ([]zkm.ZKServerInfo, error) {
	this.RLock()
	defer this.RUnlock()

	sl, ok := this.servers[class]
	if !ok {
		return []zkm.ZKServerInfo{}, fmt.Errorf("[serverlist] GetServerList failed! have no servers")
	}

	return sl, nil
}

func (this *serverlist) GetAllServerList() string {
	this.RLock()
	defer this.RUnlock()

	all := "["

	n := 0
	for _, serverinfos := range this.servers {
		if n != 0 {
			all = all + ","
		}

		s, err := json.Marshal(serverinfos)
		if err != nil {
			return err.Error()
		}

		all = all + string(s)
		n += 1
	}

	all = all + "]"

	return all
}
