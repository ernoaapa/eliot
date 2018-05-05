#!/bin/sh

set -eu

apt-get update && apt-get install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    software-properties-common

curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -
apt-key fingerprint 0EBFCD88
add-apt-repository \
   "deb https://download.docker.com/linux/debian \
   $(lsb_release -cs) \
   stable"

apt-get update && apt-get install -y docker-ce

echo "Installed docker:"
docker version