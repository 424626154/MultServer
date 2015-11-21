package main

import (
	"flag"
	"fmt"
	"lanstonetech.com/server/GateServer"
	"lanstonetech.com/server/LoginServer"
)

var ServerName *string = flag.String("name", "ServerName", "-name=ServerName")
var Index *int = flag.Int("index", 0, "-index=0")

var ServerList map[string]func(int)

func init() {
	ServerList = make(map[string]func(int))
	ServerList["GateServer"] = GateServer.Run
	ServerList["LoginServer"] = LoginServer.Run
}

func main() {
	fmt.Printf("Usage: -name ServerName(string) -index=ServerIndex(int)\n")
	flag.Parse()

	server, ok := ServerList[*ServerName]
	if !ok {
		return
	}

	server(*Index)
}
