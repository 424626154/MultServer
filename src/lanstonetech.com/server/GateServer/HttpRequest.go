package GateServer

import (
	"fmt"
	"lanstonetech.com/system/config"
	"lanstonetech.com/system/serverlist"
	"net/http"
	"time"
)

var mux map[string]func(http.ResponseWriter, *http.Request)

type myhandler struct {
}

func RunHttpServer() {
	server := http.Server{
		Addr:        fmt.Sprintf(":%d", config.SERVER_PORT),
		Handler:     &myhandler{},
		ReadTimeout: 10 * time.Second,
	}

	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	mux["/server/list"] = ServerList

	fmt.Println(server.ListenAndServe())
}

func (this myhandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h, ok := mux[r.URL.Path]
	if ok {
		h(w, r)
		return
	}

	http.NotFound(w, r)
}

func ServerList(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, serverlist.ServerList.GetAllServerList())
}
