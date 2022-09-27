#!/usr/bin/env bash
set -e
sudo iptables -t filter -F
sudo iptables -t filter -X
systemctl restart docker
