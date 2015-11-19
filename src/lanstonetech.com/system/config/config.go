package config

import (
	"fmt"
	"lanstonetech.com/common/config"
)

var (
	SERVER_TYPE  uint8
	SERVER_IP    string
	SERVER_PORT  uint16
	SERVER_GROUP string
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

func LoadServerInfo(ServerType uint8) {
	section := "LoginServer"
	index := 0

	if ServerType == 1 {
		section = "LoginServer"
	}

	server_ip := ServerConfig.MustValue(section, fmt.Sprintf("IP_%d", index))
	server_port := uint16(ServerConfig.MustInt(section, fmt.Sprintf("PORT_%d", index)))
	group := ServerConfig.MustValue(section, fmt.Sprintf("GROUP_%d", index))

	SERVER_TYPE = ServerType
	SERVER_IP = server_ip
	SERVER_PORT = server_port
	SERVER_GROUP = group
}
