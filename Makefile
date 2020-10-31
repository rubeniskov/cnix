SHELL = /bin/bash 

## Project vars
# Get name from the local directory name
PROJECT_NAME ?= $(shell basename $$(pwd))
# Get latest tag to define the version 
PROJECT_VERSION ?= $(shell git describe --abbrev=0 --tag|sed -e 's/^v//')
PROJECT_DIST_DIR ?= dist
# Get Version control reference from git
PROJECT_VC_REF = $(shell git rev-parse --short HEAD)
# Get the repo path for tracking the module directory insise GO_SRC
PROJECT_GIT_REPO_PATH = $(shell git remote -v | grep fetch | awk '{print $$2}'| sed -e 's/^git@//' | sed -e 's/.git$$//' | sed -e 's/:/\//')

# First column GOOS rest GOARCH
define PROJECT_TARGETS
darwin 386 amd64,
linux 386 amd64 arm arm64,
windows 386 amd64
endef

define PROJECT_VERSION_ERROR
PROJECT_VERSION is not set, this could be produced because there's no
available tags in your local git repository. 
  If you want to set manually use 'make $(MAKECMDGOALS) PROJECT_VERSION=v0.0.1'
endef



## Golang vars
GO_BIN ?= $(shell which go)
GO_VERSION ?= $(shell cat go.mod |sed -ne 's/^go \([0-9]\)/\1/p')
# GO_VERSION ?= $(shell $(GO_BIN) version|awk '{print $$3}'|sed -e 's/^go//')
# Get GO_SRC if exist in host machine or create a volumen to cache the dependencies
GO_SRC = $(shell [ -d $$GOPATH/src ] && echo $$GOPATH/src)
GO_IMAGE = golang:$(GO_VERSION)
# ref: https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63
GO_FLAGS = -v
# ref: https://golang.org/cmd/link/
GO_LDFLAGS = -X 'main.BuildTime=$(shell date +%s)' \
		   	 -X 'main.BuildVcRef=$(PROJECT_VC_REF)' \
		   	 -X 'main.BuildVersion=$(PROJECT_VERSION)' \
			 -X 'main.BuildSource=docker'
GO_TARGETS = $(shell awk \
				-v s="$(PROJECT_DIST_DIR)/$(PROJECT_NAME)-$(PROJECT_VERSION)" \
				'BEGIN { RS="," }{ n=split($$0,cs," ");for (i=2;i <= n;i++)printf("%s-%s-%s\n", s, cs[1],cs[i]); }'\
			<<< "$(PROJECT_TARGETS)")

## Docker vars
DOCKER_BIN = $(shell which docker)
DOCKER_RUN_FLAGS = \
		--rm \
		--volume $(or $(GO_SRC),$(shell $(DOCKER_BIN) volume create go_src)):/go/src \
		--volume $(shell pwd):/go/src/$(PROJECT_GIT_REPO_PATH) \
		--workdir /go/src/$(PROJECT_GIT_REPO_PATH)
DOCKER_BUILD_FLAGS = \
		--build-arg GO_VERSION=$(GO_VERSION) \
		--build-arg VERSION=$(PROJECT_VERSION) \
		--build-arg VC_REF=$(PROJECT_VC_REF)

INFO_VARS = PROJECT_NAME \
			PROJECT_VERSION \
			PROJECT_DIST_DIR \
			PROJECT_VC_REF \
			PROJECT_GIT_REPO_PATH \
			GO_VERSION \
			GO_BIN \
			GO_TARGETS \
			GO_SRC \
			DOCKER_BIN \
			DOCKER_RUN_FLAGS \
			DOCKER_BUILD_FLAGS


$(GO_TARGETS): check-version
	@echo "Building $@"
	$(DOCKER_BIN) run \
		$(DOCKER_RUN_FLAGS) \
		-e GO111MODULE=on -e CGO_ENABLED=0 \
		$(shell awk -F'-' '{print "-e  GOOS="$$3,"-e  GOARCH="$$4}' <<< "$@") \
		$(GO_IMAGE) \
			go build $(GO_FLAGS) \
				-ldflags="$(GO_LDFLAGS)" \
				-o $@ . 2>/dev/null

dist: install-dist $(GO_TARGETS)
	@file $(PROJECT_DIST_DIR)/*

build: install
	@echo "Building $(PROJECT_NAME)@$(PROJECT_VERSION)"
	@$(GO_BIN) build $(GO_FLAGS) \
	  	-ldflags="$(GO_LDFLAGS) -X 'main.BuildVersion=dirty' -X 'main.BuildSource=host'" \
	  	-o $(PROJECT_NAME) . >/dev/null

build-docker: check-version
	@echo "Building docker image $(PROJECT_NAME):$(PROJECT_VERSION)"
	@$(DOCKER_BIN) build $(DOCKER_BUILD_FLAGS) \
		-t $(PROJECT_NAME):$(PROJECT_VERSION) -t $(PROJECT_NAME):latest . > /dev/null
	@echo "Testing docker image $(PROJECT_NAME):$(PROJECT_VERSION)" 
	@$(DOCKER_BIN) run -it $(PROJECT_NAME):$(PROJECT_VERSION)|sed -e 's/^/ -> /'

install-dist: 
	@$(DOCKER_BIN) run $(DOCKER_RUN_FLAGS) $(GO_IMAGE) go get

install:
	@$(GO_BIN) get

check-version:
	$(if $(PROJECT_VERSION),, \
		$(error $(PROJECT_VERSION_ERROR)) \
	)

info:
	$(foreach VAR,$(INFO_VARS),\
		$(shell echo -ne >&2 "$(VAR)\n$$(echo "$($(VAR))"|sed -e 's/^/ -> /')\n"))	

clean: 
	rm -rf $(PROJECT_DIST_DIR)/*

all: info build build-dist build-docker

.PHONY: all info dist build install-dist install build-docker clean