#!/bin/bash
# 
# Push multi arch Docker image manifest
# 
set -euo pipefail

IMAGE="ernoaapa/eliotd"
PLATFORMS="linux/amd64,linux/arm64"
VERSION=$1

echo "Push multi arch manifest for version: ${VERSION} to ${IMAGE}"
echo "Platforms: ${PLATFORMS}"

manifest-tool \
  --username ${DOCKER_USER} \
  --password ${DOCKER_PASS} \
  push from-args \
  --platforms $PLATFORMS \
  --template ${IMAGE}:${VERSION}-ARCH \
  --target ${IMAGE}:${VERSION}
