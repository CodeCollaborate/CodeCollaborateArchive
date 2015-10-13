#!/usr/bin/env sh
git pull
go build server.go
sudo cp -f etc/init.d/CodeCollaborate /etc/init.d/
sudo mkdir -p /CodeCollaborate

sudo service CodeCollaborate stop
sudo cp -f server /CodeCollaborate
sudo mv -f /CodeCollaborate/server /CodeCollaborate/CodeCollaborate
sudo service CodeCollaborate start
