SHELL = /bin/bash

GO := $(shell type -P go)
GOLINT := $(shell type -P golangci-lint)

REPO := $(shell $(GO) list -m)

APP_VERSION := 0.1-alpha

GO_BUILD_ARGS := 
# GO_BUILD_ARGS := \
#   -ldflags " \
#     -X 'main.Version=$(APP_VERSION)' \
# "
#   -race \
#     -X '$(REPO)/internal/version.Version=$(APP_VERSION)' \


export GO111MODULE = on
export CGO_ENABLED = 0

BIN_SOURCES := $(wildcard cmd/*/)
BIN_NAMES   := $(patsubst cmd/%/,%,$(BIN_SOURCES))
BIN_DESTS   := $(addprefix bin/,$(BIN_NAMES))

.PHONY: build
build: fix fetch clean $(BIN_DESTS)
	# $(MAKE) $@

.PHONY: clean
clean: bin
	# $(MAKE) $@
	rm -vf $(BIN_DESTS)

.PHONY: fix
fix:
	# $(MAKE) $@
	$(GO) mod tidy
	$(GO) fmt -mod=mod ./...

.PHONY: fetch
fetch: fix
	# $(MAKE) $@
	$(GO) mod download -x
	$(GO) mod vendor -v

bin:
	# $(MAKE) $@
	mkdir bin

bin/%: bin
	# $(MAKE) $@
	$(GO) build \
	  -o $@ \
	  $(GO_BUILD_ARGS) \
	  $(PWD)/cmd/$(notdir $@)

.PHONY: test
test:
	$(GO) test -v ./...

.PHONY: lint
lint:
	$(GOLINT) run ./... \
	  --verbose \
	  --skip-dirs-use-default \
	  --allow-parallel-runners

.PHONY: lint-fix
lint-fix:
	$(GOLINT) run ./... \
	  --fix \
	  --verbose \
	  --skip-dirs-use-default \
	  --allow-parallel-runners
