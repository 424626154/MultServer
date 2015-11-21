#!/bin/bash

rm -f ./server
CGO_ENABLED=0 GOOS=linux GOARCH=amd64  go build -v ./server.go

if [ -a server ] ; then
	echo "==>build Server successed!"
else
	echo "==>build Server failed!"
	exit -1
fi
