#!/bin/sh

set -eu

echo "GOPATH: $GOPATH (PATH: $PATH)"

go get -u github.com/goreleaser/goreleaser
goreleaser --version

go get -u github.com/goreleaser/nfpm
cd $GOPATH/src/github.com/goreleaser/nfpm
go get ./...
GOBIN=$GOPATH/bin go install ./cmd/nfpm/
cd -
nfpm --version

go get -u github.com/estesp/manifest-tool
manifest-tool --version

go get -u github.com/mlafeldt/pkgcloud/...