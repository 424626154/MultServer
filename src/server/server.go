package main

import (
	"flag"
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
	// if len(os.Args) < 2 {
	// 	fmt.Printf("ERR: Less Argument!\nUsage: %s LoginServer\n", os.Args[0])
	// 	return
	// }
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		return
	}

	server, ok := ServerList[args[0]]
	if !ok {
		return
	}
}
