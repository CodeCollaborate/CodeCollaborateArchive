#!/usr/bin/env sh
echo "Installing dependencies"
go get github.com/gorilla/websocket
go get gopkg.in/mgo.v2
go get golang.org/x/crypto/bcrypt