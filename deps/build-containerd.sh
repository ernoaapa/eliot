#!/bin/sh
# 
# 

set -eu

CONTAINERD_VERSION="v1.1.3"
RUNC_DEP_VERSION="1.0.0-rc4+a618ab5a"
BUILD_DIR=/tmp/build
ELIOT_SRC_DIR=$(pwd)
TARGET_DIR=${ELIOT_SRC_DIR}/dist
GOOS=${GOOS-"linux"}
GOARCH=${GOARCH} # required
GOARM=${GOARM:-""} # optional

mkdir -p $TARGET_DIR

rm -rf ${BUILD_DIR} && mkdir -p ${BUILD_DIR}/src/github.com/containerd/containerd
wget -qO- "https://github.com/containerd/containerd/archive/${CONTAINERD_VERSION}.tar.gz" | tar xvz --strip-components=1 -C ${BUILD_DIR}/src/github.com/containerd/containerd
cd ${BUILD_DIR}/src/github.com/containerd/containerd
# Modify Makefile to have optional VERSION and REVISION
sed -i='' 's/VERSION=/VERSION?=/g' Makefile
sed -i='' 's/REVISION=/REVISION?=/g' Makefile

echo "Compile containerd os:${GOOS}, arch:${GOARCH}, variant:${GOARM}"
GOPATH=${BUILD_DIR} VERSION=${CONTAINERD_VERSION} REVISION=unknown make binaries

cat << EOF > ./nfpm.yaml
name: "containerd"
arch: "${GOARCH}${GOARM:-""}"
platform: "${GOOS}"
version: "${CONTAINERD_VERSION}"
section: "default"
priority: "extra"
depends:
- runc (>=${RUNC_DEP_VERSION})
maintainer: "Erno Aapa <erno.aapa@gmail.com>"
description: An open and reliable container runtime
homepage: "https://containerd.io"
license: "Apache 2.0"
bindir: "/usr/local/bin"
files:
  ./bin/containerd: "/usr/local/bin/containerd"
  ./bin/ctr: "/usr/local/bin/ctr"
  ./bin/containerd-shim: "/usr/local/bin/containerd-shim"
  ${ELIOT_SRC_DIR}/build/etc/systemd/system/containerd.service: "/etc/systemd/system/containerd.service"
EOF

nfpm pkg --target ${TARGET_DIR}/containerd_${CONTAINERD_VERSION}_${GOOS}_${GOARCH}${GOARM:-""}.deb
