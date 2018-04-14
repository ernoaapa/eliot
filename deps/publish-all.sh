#!/bin/sh

set -eu

BASEDIR=$(dirname "$0")
if [ -z ${PACKAGECLOUD_TOKEN+x} ]; then
  echo "You must define PACKAGECLOUD_TOKEN environment variable to publish packages"
  exit 1
fi

rm -rf ./dist
${BASEDIR}/build-containerd.sh
${BASEDIR}/build-runc.sh

for package in dist/*.deb; do
  pkgcloud-push ernoaapa/eliot/raspbian/stretch $package || echo "${package} upload failed. Already uploaded?"
done