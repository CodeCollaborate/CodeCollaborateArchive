#!/usr/bin/env sh
git pull
go build CodeCollaborate.go
sudo cp -f etc/init.d/CodeCollaborate /etc/init.d/
sudo mkdir -p /CodeCollaborate

sudo service CodeCollaborate stop
sudo cp -f CodeCollaborate /CodeCollaborate
sudo service CodeCollaborate start
