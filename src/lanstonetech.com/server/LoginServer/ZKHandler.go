package LoginServer

import (
	"lanstonetech.com/common/logger"
	"lanstonetech.com/system/config"
	"lanstonetech.com/system/serverlist"
	"lanstonetech.com/system/zkm"
	"time"
)

func InitZK() {
	zkm.Init()

	serverlist.ServerList.AddMonitor(config.SERVER_GROUP, config.SERVER_TYPE)
	zkm.AddObservers(&serverlist.ServerList)

	zkm.Start()

	go registerZK()
}

func registerZK() {
	for {
		err := registerZKNetwork()
		if err != nil {
			time.Sleep(5 * time.Second)
			// logger.Errorf("registerZKNetwork failed! err=%v", err)
			continue
		}

		time.Sleep(5 * time.Second)
	}
}

func registerZKNetwork() error {
	defer logger.CatchException()

	err := zkm.Server.Register(config.SERVER_GROUP, config.SERVER_TYPE, 0, config.SERVER_IP, config.SERVER_PORT, config.SERVER_IP, config.SERVER_PORT+100, "www.lanstonetech.com:8080")
	if err != nil {
		return err
	}

	serversinfo, err := serverlist.ServerList.GetServerList(config.SERVER_TYPE)
	if err != nil {
		return err
	}

	for i, serverinfo := range serversinfo {
		logger.Errorf("serverinfo[%d] = %#v", i, serverinfo)
	}

	return nil
}
