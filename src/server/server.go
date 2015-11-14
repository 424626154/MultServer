package main

import (
	"fmt"
	"lanstonetech.com/server/LoginServer"
	"os"
)

var ServerList map[string]func()

func init() {
	ServerList = make(map[string]func())
	ServerList["LoginServer"] = LoginServer.Run
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("ERR: Less Argument!\nUsage: %s LoginServer\n", os.Args[0])
		return
	}

	ServerList[os.Args[1]]()
}
