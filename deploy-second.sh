#!/usr/bin/env bash
set -e
rm -r /home/go-user-api/
mkdir  /home/go-user-api/
cd /home/go-user-api/
git clone git@bitbucket.org:iliaposmac/trucktrace-user-api.git .
random=$RANDOM
docker stop golang-user
docker build -t golang-user:${random} .
docker run --name golang-user --network host -d --rm golang-user:${random}
