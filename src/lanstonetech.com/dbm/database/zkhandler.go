package database

import (
	"encoding/json"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"lanstonetech.com/common/logger"
	"lanstonetech.com/system/zkm"
	"sync"
	"time"
)

var (
	DB_DRIVER_POSTGRES = "postgres"
)

type Observer interface {
	Group() string
	OnUpdateDatabaseList(master, slaves []zkm.DBServerInfo) error
}

type SQLServer struct {
	sync.RWMutex
	obsevers []Observer
	master   []zkm.DBServerInfo
	slaves   []zkm.DBServerInfo
}

func (this *SQLServer) AddObserver(observer Observer) {
	this.obsevers = append(this.obsevers, observer)

	zkm.AddObserver(zkm.MakeDatabaseGroup(DB_DRIVER_POSTGRES, observer.Group()), this)

	//request when register
	this.requestDatabaseList(observer)
}

func (this *SQLServer) Nodes() ([]string, error) {
	return []string{}, nil
}

func (this *SQLServer) Request() error {
	this.RequestDabaseList()
}

func (this *SQLServer) RequestDabaseList() {
	for _, observer := range this.obsevers {
		go this.requestDabaseList(observer)
	}
}

func (this *SQLServer) requestDatabaseList(observer Observer) {
	defer logger.CatchException()

	for {
		sl, _, err := zkm.ChildrenW(zkm.MakeDatabaseGroup(DB_DRIVER_POSTGRES, observer.Group()))
		if err != nil {
			logger.Errorf("zkm.ChildrenW failed! err=%v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		var master []zkm.DBServerInfo
		var slaves []zkm.DBServerInfo
		var node zkm.DBServerInfo

		for _, s := range sl {
			if err := json.Unmarshal([]byte(s), &node); err != nil {
				logger.Errorf("requestDatabaseList failed! parse.err=%v", err)
				continue
			}

			if len(node.Slaveof) == 0 {
				master = append(master, node)
			} else {
				slaves = append(slaves, node)
			}
		}

		this.Lock()
		defer this.Unlock()
		this.master = master
		this.slaves = slaves

		observer.OnUpdateDatabaseList(master, slaves)
		return
	}
}

func (this *SQLServer) ProcessException(event *zkm.EventArgs) error {
	return this.Request()
}

func (this *SQLServer) ProcessEvent(event *zkm.EventArgs) error {
	if event.Event.Type == zk.EventNodeCreated {

	} else if event.Event.Type == zk.EventNodeDeleted {

	} else if event.Event.Type == zk.EventNodeDataChanged {

	} else if event.Event.Type == zk.EventNodeChildrenChanged {
		//子节点数据改变
		this.RequestDabaseList()
	}
	return nil
}
