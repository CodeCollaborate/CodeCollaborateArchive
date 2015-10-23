#!/usr/bin/env sh
echo "Building binary from go project"
go build CodeCollaborate.go

echo "Adding init.d entry"
sudo cp -f etc/init.d/CodeCollaborate /etc/init.d/
sudo mkdir -p /CodeCollaborate

echo "Restarting and updating service"
sudo service CodeCollaborate stop
sudo cp -f CodeCollaborate /CodeCollaborate
sudo service CodeCollaborate start
