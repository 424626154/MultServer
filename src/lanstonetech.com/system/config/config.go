package config

import (
	"fmt"
	"lanstonetech.com/common/config"
)

var (
	// SERVER       string
	SERVER_IP    string
	SERVER_PORT  string
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

func LoadServerInfo(ServerType int) {
	section := "LoginServer"
	index := 0

	if ServerType == 1 {
		section = "LoginServer"
	}

	server_ip := ServerConfig.MustValue(section, fmt.Sprintf("IP_%d", index))
	server_port := ServerConfig.MustValue(section, fmt.Sprintf("PORT_%d", index))
	group := ServerConfig.MustValue(section, fmt.Sprintf("GROUP_%d", index))

	SERVER_IP = server_ip
	SERVER_PORT = server_port
	SERVER_GROUP = group
}
