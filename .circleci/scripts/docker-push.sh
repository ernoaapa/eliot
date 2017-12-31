#!/bin/bash -e

REGISTRY=ernoaapa
PLATFORMS="linux/amd64,linux/arm64"
GIT_HASH=$(git describe --tags --always --dirty)


for bin in eliotd; do
  image="${REGISTRY}/${bin}"
  
  for osarch in ${PLATFORMS//,/ }; do
    os="${osarch%%/*}"
    arch="${osarch#*/}"
    version="${GIT_HASH}-${arch}"
    tag="${image}:${version}"
    echo "Building container for: $bin $os $arch, tag: ${tag}"

    docker push ${tag}
  done

	manifest-tool \
		--username ${DOCKER_USER} \
		--password ${DOCKER_PASS} \
		push from-args \
    --platforms $PLATFORMS \
    --template ${image}:${GIT_HASH}-ARCH \
    --target ${image}:${GIT_HASH}
done
