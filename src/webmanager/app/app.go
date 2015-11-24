package app

import (
	"log"
	"net/http"
	"time"
	"webmanager/routers"
)

var Session *http.Server

func Run() {
	Session = &http.Server{
		Addr:           ":8888",
		Handler:        &routers.Dax,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(Session.ListenAndServe())
}
