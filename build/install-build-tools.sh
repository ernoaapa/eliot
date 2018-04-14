#!/bin/sh

set -eu

go get -u github.com/goreleaser/goreleaser
# Switch to fixed fork
cd $GOPATH/src/github.com/goreleaser/goreleaser
git remote add ernoaapa https://github.com/ernoaapa/goreleaser.git || true
git fetch ernoaapa
git checkout fix-arm-architecture
GOBIN=$GOPATH/bin go install .

go get -u github.com/goreleaser/nfpm
# Switch to fixed fork
cd $GOPATH/src/github.com/goreleaser/nfpm
git remote add ernoaapa https://github.com/ernoaapa/nfpm.git || true
git fetch ernoaapa
git checkout fix-arm-architecture
GOBIN=$GOPATH/bin go install ./cmd/nfpm/

go get github.com/estesp/manifest-tool
go get github.com/mlafeldt/pkgcloud/...