#!/bin/sh
# 
# 

set -eu

RUNC_VERSION="v0.1.1"
BUILD_DIR=/tmp/build
TARGET_DIR=$(pwd)/dist
GOOS=${GOOS-"linux"}
GOARCH=${GOARCH} # required
GOARM=${GOARM:-""} # optional

mkdir -p $TARGET_DIR

rm -rf ${BUILD_DIR} && mkdir -p ${BUILD_DIR}/src/github.com/opencontainers/runc
wget -qO- "https://github.com/opencontainers/runc/archive/${RUNC_VERSION}.tar.gz" | tar xvz --strip-components=1 -C ${BUILD_DIR}/src/github.com/opencontainers/runc
cd ${BUILD_DIR}/src/github.com/opencontainers/runc

echo "Build runc os:${GOOS}, arch:${GOARCH}, variant:${GOARM}"
GOPATH=${BUILD_DIR} make

cat << EOF > ./nfpm.yaml
name: "runc"
arch: "${GOARCH}${GOARM:-""}"
platform: "${GOOS}"
version: "${RUNC_VERSION}"
section: "default"
priority: "extra"
depends:
- libseccomp2
maintainer: "Erno Aapa <erno.aapa@gmail.com>"
description: runc is a CLI tool for spawning and running containers according to the OCI specification.
homepage: "https://www.opencontainers.org"
license: "Apache 2.0"
bindir: "/usr/local/bin"
files:
  ./runc: "/usr/local/bin/runc"
EOF

nfpm pkg --target ${TARGET_DIR}/runc_${RUNC_VERSION}_${GOOS}_${GOARCH}${GOARM:-""}.deb
