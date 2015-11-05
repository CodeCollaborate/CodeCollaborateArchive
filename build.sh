#!/usr/bin/env sh
echo "Building server binary from go project"
go build CodeCollaborate.go

echo "Building scrunching jar from source"
cd Scrunching && mvn package && cd ..

echo "Adding init.d entry"
sudo cp -f etc/init.d/CodeCollaborate /etc/init.d/
sudo mkdir -p /CodeCollaborate

echo "Restarting and updating service"
sudo service CodeCollaborate stop
sudo cp -f CodeCollaborate /CodeCollaborate
sudo cp -f Scrunching/target/Scrunching.jar /CodeCollaborate
sudo service CodeCollaborate start
