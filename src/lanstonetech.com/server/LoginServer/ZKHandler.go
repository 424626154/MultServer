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

	go registerZK()
}

func registerZK() {
	for {
		err := registerZKNetwork()
		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}

		time.Sleep(15 * time.Second)
	}
}

func registerZKNetwork() error {
	defer logger.CatchException()

	err := zkm.Server.Register(config.SERVER_GROUP, config.SERVER_TYPE, 0, config.SERVER_IP, config.SERVER_PORT, "", config.SERVER_PORT+100, "")
	if err != nil {
		return err
	}

	logger.Errorf("registerZKNetwork successful!")
	return nil
}
