package GateServer

import (
	"lanstonetech.com/common"
	"lanstonetech.com/common/logger"
	"lanstonetech.com/system/config"
	"lanstonetech.com/system/serverlist"
	"lanstonetech.com/system/zkm"
	"time"
)

func InitZK() {
	zkm.Init()

	serverlist.ServerList.AddMonitor(config.SERVER_GROUP, common.GATE_SERVER_TYPE)
	serverlist.ServerList.AddMonitor(config.SERVER_GROUP, common.LOGIN_SERVER_TYPE)
	zkm.AddObservers(&serverlist.ServerList)

	zkm.Start()

	go registerZK()
}

func registerZK() {
	for {
		err := registerZKNetwork()
		if err != nil {
			time.Sleep(15 * time.Second)
			logger.Errorf("registerZKNetwork failed! err=%v", err)
			continue
		}

		time.Sleep(15 * time.Second)
	}
}

func registerZKNetwork() error {
	defer logger.CatchException()

	err := zkm.Server.Register(config.SERVER_GROUP, config.SERVER_TYPE, config.SERVER_INDEX, config.SERVER_IP, config.SERVER_PORT, "", 0, "www.lanstonetech.com:8080")
	if err != nil {
		return err
	}

	return nil
}
