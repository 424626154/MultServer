package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pg"
	"math/rand"
	"time"
)

type ServerInfo struct {
	Group    string
	IP       string
	Port     int
	DBName   string
	User     string
	Password string
	Slave    int
	Slaveof  string
	Offline  bool
}

type ServerGroup struct {
	Master ServerInfo
	Slaves []ServerInfo
}

type ServerPool struct {
	Items       []ServerGroup
	connections map[string][]*sql.DB
}

func (this *ServerPool) init() error {
	for _, item := range this.Items {
		conn, err := this.connectDB(item.Master.IP, item.Master.Port, item.Master.DBName, item.Master.User, item.Master.Password)
		if err != nil {
			continue
		}

		var conns []*sql.DB
		conns = append(conns, conn)

		for _, slave := range item.Slaves {
			if slave.Offline {
				continue
			}

			conn, err := this.connectDB(slave.IP, slave.Port, slave.DBName, slave.User, slave.Password)
			if err != nil {
				continue
			}

			conns = append(conns, conn)
		}

		this.connections[fmt.Sprintf("%s:%s:%s", item.Master.IP, item.Master.Port, item.Master.DBName)] = conns
	}

	if len(this.connections) == 0 {
		return fmt.Errorf("ServerPool init failed! connections==nil")
	}

	return nil
}

func (this *ServerPool) AddMaster(server ServerInfo) error {
	var item ServerGroup
	item.Master = server
	this.Items = append(this.Items, item)

	return nil
}

func (this *ServerPool) AddSlave(server ServerInfo) error {
	for i, _ := range this.Items {
		name := fmt.Sprintf("%s:%s:%s", this.Items[i].Master.IP, this.Items[i].Master.Port, this.Items[i].Master.DBName)
		if server.Slaveof == name {
			this.Items[i].Slaves == append(this.Items[i].Slave, server)
			return nil
		}
	}

	return fmt.Errorf("AddSlave failed! No master!")
}

func (this *ServerPool) release() error {
	for _, item := range this.connections {
		for _, conn := range item {
			conn.Close()
		}
	}

	this.connections = nil
}

func (this *ServerPool) getconn(isWrite bool) (*sql.DB, error) {
	if len(this.connections) == 0 {
		return nil, fmt.Errorf("getconn failed! master is empty")
	}

	var dbs []*sql.DB
	for _, item := range this.connections {
		if len(item) > 0 {
			dbs = item
			break
		}
	}

	count := len(dbs)
	if count == 1 {
		return dbs[0], nil
	} else if isWrite {
		return dbs[0], nil
	} else {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		index := rand.Int31n(count-1) + 1
		return dbs[int(index)], nil
	}
}

func (ServerPool) connectDB(IP string, Port int, dbname, user, pass string) (*sql.DB, error) {

	/*
	 * dbname - The name of the database to connect to
	 * user - The user to sign in as
	 * password - The user's password
	 * host - The host to connect to. Values that start with / are for unix domain sockets. (default is localhost)
	 * port - The port to bind to. (default is 5432)
	 * sslmode - Whether or not to use SSL (default is require, this is not the default for libpq)
	 * fallback_application_name - An application_name to fall back to if one isn't provided.
	 * connect_timeout - Maximum wait for connection, in seconds. Zero or not specified means wait indefinitely.
	 * sslcert - Cert file location. The file must contain PEM encoded data.
	 * sslkey - Key file location. The file must contain PEM encoded data.
	 * sslrootcert - The location of the root certificate file. The file must contain PEM encoded data.
	 */

	session, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable", user, pass, IP, Port, dbname))
	if err != nil {
		return nil, err
	}

	return session, nil
}
