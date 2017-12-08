#!/bin/bash -e

REGISTRY=ernoaapa
PLATFORMS="linux/amd64,linux/arm64"
GIT_HASH=$(git describe --tags --always --dirty)


for bin in eliot-discovery eliotd; do
  image="${REGISTRY}/${bin}"
  
  for osarch in ${PLATFORMS//,/ }; do
    os="${osarch%%/*}"
    arch="${osarch#*/}"
    echo "Building container for: $bin $os $arch"

    sed \
	    -e "s|ARG_BIN|${bin}|g" \
			-e "s|ARG_OS|${os}|g" \
			-e "s|ARG_ARCH|${arch}|g" \
	    Dockerfile.tmpl > .dockerfile-${arch}

    version="${GIT_HASH}-${arch}"
    tag="${image}:${version}"
    echo "Build docker image ${tag}:"
    echo "---------------Dockerfile start---------------------"
    cat .dockerfile-${arch}
    echo "----------------Dockerfile end----------------------"

    docker build -t ${tag} -f .dockerfile-${arch} .
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
