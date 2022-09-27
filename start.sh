#!/usr/bin/env bash
echo "Go mod inited"
screen -dmS golang-user go run cmd/app/*.go
echo "Finished"
