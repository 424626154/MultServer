package database

import (
	"lanstonetech.com/common/logger"
	"lanstonetech.com/system/zkm"
	"sync"
)

type SQLConn struct {
	sync.RWMutex
	Group string
	pool  ServerPool
}

func (this *SQLConn) Init(group string) {
	this.Group = group

	SQLServer.AddObserver(this)
}

func (this *SQLConn) Group() string {
	return this.Group()
}

func (this *SQLConn) OnUpdateDatabaseList(master, slaves []zkm.DBServerInfo) error {
	var pool ServerPool
	var server ServerInfo

	for _, item := range master {
		if item.Group != this.Group() {
			continue
		}

		server.Group = item.Group
		server.IP = item.IP
		server.Port = item.Port
		server.DBName = item.DBName
		server.User = item.User
		server.Password = item.Password
		server.Slave = 0
		if item.Offline == 1 {
			server.Offline = true
		} else {
			server.Offline = false
		}

		pool.AddMaster(server)
	}

	for _, item := range slaves {
		if item.Group != this.Group() {
			continue
		}

		server.Group = item.Group
		server.IP = item.IP
		server.Port = item.Port
		server.DBName = item.DBName
		server.User = item.User
		server.Password = item.Password
		server.Slave = 1
		if item.Offline == 1 {
			server.Offline = true
		} else {
			server.Offline = false
		}

		pool.AddSlave(server)
	}

	return this.updatePool(&pool)
}

func (this *SQLConn) updatePool(pool *ServerPool) error {
	if err := pool.init(); err != nil {
		return err
	}

	this.Lock()
	defer this.Unlock()

	this.pool.release()
	this.pool = *pool
	return nil
}
