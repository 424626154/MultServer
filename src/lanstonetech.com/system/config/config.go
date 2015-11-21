package config

import (
	"fmt"
	"lanstonetech.com/common"
	"lanstonetech.com/common/config"
)

var (
	SERVER_TYPE  uint8
	SERVER_IP    string
	SERVER_PORT  uint16
	SERVER_GROUP string
	SERVER_INDEX int
)

var ServerConfig config.ConfigFile

func init() {
	ini := "../../conf/server.ini"
	config, err := config.LoadConfigFile(ini)
	if err != nil {
		panic(err)
	}

	ServerConfig = *config
}

func LoadServerInfo(ServerType, index int) {
	var section string
	if ServerType == common.LOGIN_SERVER_TYPE {
		section = "LoginServer"

	} else if ServerType == common.GATE_SERVER_TYPE {
		section = "GateServer"
	} else {
		return
	}

	server_ip := ServerConfig.MustValue(section, fmt.Sprintf("IP_%d", index))
	server_port := uint16(ServerConfig.MustInt(section, fmt.Sprintf("PORT_%d", index)))
	group := ServerConfig.MustValue(section, fmt.Sprintf("GROUP_%d", index))

	SERVER_TYPE = uint8(ServerType)
	SERVER_IP = server_ip
	SERVER_PORT = server_port
	SERVER_GROUP = group
	SERVER_INDEX = index
}
