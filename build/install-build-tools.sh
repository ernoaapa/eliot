#!/bin/sh

set -eu

echo "GOPATH: $GOPATH"

rm -rf $GOPATH/src/github.com/goreleaser/nfpm

go get -u github.com/goreleaser/goreleaser

go get -u github.com/goreleaser/nfpm
# Switch to fixed fork
cd $GOPATH/src/github.com/goreleaser/nfpm
git remote add ernoaapa https://github.com/ernoaapa/nfpm.git || true
git fetch ernoaapa
git checkout ernoaapa/fix-arm-architecture
go get ./...
GOBIN=$GOPATH/bin go install ./cmd/nfpm/

go get github.com/estesp/manifest-tool
go get github.com/mlafeldt/pkgcloud/...