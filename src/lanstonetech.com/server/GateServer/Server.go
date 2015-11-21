package GateServer

import (
	"lanstonetech.com/common"
	"lanstonetech.com/common/logger"
	"lanstonetech.com/system/config"
)

func InitConf(index int) {
	config.LoadServerInfo(common.GATE_SERVER_TYPE, index)
	logger.Infof("ip=%v port=%v group=%v", config.SERVER_IP, config.SERVER_PORT, config.SERVER_GROUP)
}

func InitLog() {
	logger.SetConsole(true)
	logger.SetConsolePrefix("GateServer")
	// logger.Initialize("./log", "LoginServer_0")
	logger.SetLevel(logger.LEVEL(2))
}

func Run(index int) {
	defer logger.CatchException()

	logger.Errorf("index=%v", index)
	InitLog()
	InitConf(index)
	logger.Infof("server start...")

	InitZK()

	RunHttpServer()
}
