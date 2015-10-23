#!/usr/bin/env sh
echo "Pulling latest version"
git pull
./install-dependencies.sh
./build.sh