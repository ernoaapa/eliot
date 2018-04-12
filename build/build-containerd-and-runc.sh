#!/bin/sh
# 
# 

set -eu

CONTAINERD_VERSION="v1.0.3"
RUNC_VERSION="v0.1.1"
BUILD_DIR=/tmp/build
TARGET_DIR=$(pwd)/bin

mkdir -p $TARGET_DIR

# CONTAINERD
echo "Build containerd"
rm -rf ${BUILD_DIR} && mkdir -p ${BUILD_DIR}/src/github.com/containerd/containerd
wget -qO- "https://github.com/containerd/containerd/archive/${CONTAINERD_VERSION}.tar.gz" | tar xvz --strip-components=1 -C ${BUILD_DIR}/src/github.com/containerd/containerd
cd ${BUILD_DIR}/src/github.com/containerd/containerd
# Modify Makefile to have optional VERSION and REVISION
sed -i='' 's/VERSION=/VERSION?=/g' Makefile
sed -i='' 's/REVISION=/REVISION?=/g' Makefile
GOPATH=${BUILD_DIR} VERSION=${CONTAINERD_VERSION} REVISION=unknown make binaries
cp bin/containerd ${TARGET_DIR}
cp bin/ctr ${TARGET_DIR}
cp bin/containerd-shim ${TARGET_DIR}

# RUNC
echo "Build runc"
rm -rf ${BUILD_DIR} && mkdir -p ${BUILD_DIR}/src/github.com/opencontainers/runc
wget -qO- "https://github.com/opencontainers/runc/archive/${RUNC_VERSION}.tar.gz" | tar xvz --strip-components=1 -C ${BUILD_DIR}/src/github.com/opencontainers/runc
cd ${BUILD_DIR}/src/github.com/opencontainers/runc
make
cp runc ${TARGET_DIR}
