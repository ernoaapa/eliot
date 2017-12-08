
# Modified version of
# https://github.com/thockin/go-build-template

# The binary to build (just the basename).
BIN ?= eliotd

# This repo's root import path (under GOPATH).
PKG := github.com/ernoaapa/eliot

# Where to push the docker image.
REGISTRY ?= ernoaapa

# Which architecture to build
ARCH ?= amd64
OS ?= linux

# This version-strategy uses git tags to set the version string
GIT_HASH ?= $(shell git describe --tags --always --dirty)
VERSION := $(GIT_HASH)-$(ARCH)
#
# This version-strategy uses a manual value to set the version string
#VERSION := 1.2.3

###
### These variables should not need tweaking.
###

SRC_DIRS := cmd pkg # directories which hold app source (not vendored)

ALL_PLATFORMS := darwin-amd64 linux-amd64 linux-arm64
CONTAINER_PLATFORMS := linux-amd64 linux-arm64

IMAGE := $(REGISTRY)/$(BIN)

BUILD_IMAGE ?= golang:1.9-alpine

# If you want to build all binaries, see the 'all-build' rule.
# If you want to build all containers, see the 'all-container' rule.
# If you want to build AND push all containers, see the 'all-push' rule.
all: build

build-%:
	@$(MAKE) --no-print-directory ARCH=$(word 2,$(subst -, ,$*)) OS=$(word 1,$(subst -, ,$*)) build

container-%:
	@$(MAKE) --no-print-directory ARCH=$(word 2,$(subst -, ,$*)) OS=$(word 1,$(subst -, ,$*))  container

push-%:
	@$(MAKE) --no-print-directory ARCH=$(word 2,$(subst -, ,$*)) OS=$(word 1,$(subst -, ,$*))  push

all-build: $(addprefix build-, $(ALL_PLATFORMS))

all-container: $(addprefix container-, $(CONTAINER_PLATFORMS))

all-push: $(addprefix push-, $(CONTAINER_PLATFORMS))

build: bin/$(OS)-$(ARCH)/$(BIN)

 bin/$(OS)-$(ARCH)/$(BIN): build-dirs
 ifeq ($(filter $(OS)-$(ARCH),$(ALL_PLATFORMS)),)
 	$(error unsupported platform $(OS)-$(ARCH) not in $(ALL_PLATFORMS))
 endif
	@echo "building: $@"
	@docker run                                                                   \
	    -ti                                                                       \
	    --rm                                                                      \
	    -u $$(id -u):$$(id -g)                                                    \
	    -v "$$(pwd)/.go:/go"                                                      \
	    -v "$$(pwd):/go/src/$(PKG)"                                               \
	    -v "$$(pwd)/bin/$(OS)-$(ARCH):/go/bin"                                    \
	    -v "$$(pwd)/bin/$(OS)-$(ARCH):/go/bin/$(OS)_$(ARCH)"                      \
	    -v "$$(pwd)/.go/std/$(OS)-$(ARCH):/usr/local/go/pkg/$(OS)_$(ARCH)_static" \
	    -w /go/src/$(PKG)                                                         \
	    $(BUILD_IMAGE)                                                            \
	    /bin/sh -c "                                                              \
	        ARCH=$(ARCH)                                                          \
					OS=$(OS)                                                              \
	        VERSION=$(VERSION)                                                    \
	        PKG=$(PKG)                                                            \
	        ./build/build.sh                                                      \
	    "

# Example: make shell CMD="-c 'date > datefile'"
shell: build-dirs
	@echo "launching a shell in the containerized build environment"
	@docker run                                                                   \
	    -ti                                                                       \
	    --rm                                                                      \
	    -u $$(id -u):$$(id -g)                                                    \
	    -v "$$(pwd)/.go:/go"                                                      \
	    -v "$$(pwd):/go/src/$(PKG)"                                               \
	    -v "$$(pwd)/bin/$(OS)-$(ARCH):/go/bin"                                    \
	    -v "$$(pwd)/bin/$(OS)-$(ARCH):/go/bin/$(OS)_$(ARCH)"                      \
	    -v "$$(pwd)/.go/std/$(OS)-$(ARCH):/usr/local/go/pkg/$(OS)_$(ARCH)_static" \
	    -w /go/src/$(PKG)                                                         \
	    $(BUILD_IMAGE)                                                            \
	    /bin/sh $(CMD)

DOTFILE_IMAGE = $(subst :,_,$(subst /,_,$(IMAGE))-$(VERSION))

container: .container-$(DOTFILE_IMAGE) container-name
.container-$(DOTFILE_IMAGE): bin/$(OS)-$(ARCH)/$(BIN) Dockerfile.in
	@sed \
	    -e 's|ARG_BIN|$(BIN)|g' \
			-e 's|ARG_OS|$(OS)|g' \
			-e 's|ARG_ARCH|$(ARCH)|g' \
	    Dockerfile.in > .dockerfile-$(ARCH)
	@docker build -t $(IMAGE):$(VERSION) -f .dockerfile-$(ARCH) .
	@docker images -q $(IMAGE):$(VERSION) > $@

container-name:
	@echo "container: $(IMAGE):$(VERSION)"

push: .push-$(DOTFILE_IMAGE) push-name
.push-$(DOTFILE_IMAGE): .container-$(DOTFILE_IMAGE)
ifeq ($(findstring gcr.io,$(REGISTRY)),gcr.io)
	@gcloud docker -- push $(IMAGE):$(VERSION)
else
	@docker push $(IMAGE):$(VERSION)
endif
	@docker images -q $(IMAGE):$(VERSION) > $@

push-name:
	@echo "pushed: $(IMAGE):$(VERSION)"

push-manifest:
	manifest-tool \
		--username $(DOCKER_USER) \
		--password $(DOCKER_PASS) \
		push from-args \
    --platforms linux/amd64,linux/arm64 \
    --template $(IMAGE):$(GIT_HASH)-ARCH \
    --target $(IMAGE):$(GIT_HASH)

version:
	@echo $(VERSION)

test: build-dirs
	@docker run                                                                   \
	    -ti                                                                       \
	    --rm                                                                      \
	    -u $$(id -u):$$(id -g)                                                    \
	    -v "$$(pwd)/.go:/go"                                                      \
	    -v "$$(pwd):/go/src/$(PKG)"                                               \
	    -v "$$(pwd)/bin/$(ARCH):/go/bin"                                          \
	    -v "$$(pwd)/.go/std/$(OS)-$(ARCH):/usr/local/go/pkg/$(OS)_$(ARCH)_static" \
	    -w /go/src/$(PKG)                                                         \
	    $(BUILD_IMAGE)                                                            \
	    /bin/sh -c "                                                              \
	        ./build/test.sh $(SRC_DIRS)                                           \
	    "

build-dirs:
	@mkdir -p bin/$(OS)-$(ARCH)
	@mkdir -p .go/src/$(PKG) .go/pkg .go/bin .go/std/$(OS)-$(ARCH)

clean: container-clean bin-clean

container-clean:
	rm -rf .container-* .dockerfile-* .push-*

bin-clean:
	rm -rf .go bin

publish-docs:
	./build/docs.sh