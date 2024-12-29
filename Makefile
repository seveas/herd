GOGCFLAGS := -gcflags=all=-e

# Let's not rebuild the parser if we don't have antlr available
ifeq ("", "$(strip $(shell which antlr))")
	antlr_sources :=
else
	antlr_sources := scripting/parser/herd_base_listener.go scripting/parser/herd_lexer.go scripting/parser/herd_listener.go scripting/parser/herd_parser.go
endif

# Let's not rebuild the protobuf code if we don't have protobuf available
ifeq ("", "$(strip $(shell which protoc))")
	protobuf_sources :=
else ifeq ("", "$(strip $(shell which protoc-gen-go))")
	protobuf_sources :=
else ifeq ("", "$(strip $(shell which protoc-gen-go-crpc))")
	protobuf_sources :=
else
	protobuf_sources = provider/plugin/common/plugin.pb.go provider/plugin/common/plugin_grpc.pb.go
endif

# The main program
herd: go.mod go.sum *.go cmd/herd/*.go ssh/*.go scripting/*.go provider/*/*.go provider/plugin/common/*.go $(protobuf_sources) $(antlr_sources)
	go build $(GOGCFLAGS) -o "$@" github.com/seveas/herd/cmd/herd

# External providers
provider_plugins := aws azure consul google prometheus puppet tailscale transip
cmd/herd-provider-%/main.go: cmd/herd-provider-example/main.go
	mkdir -p cmd/herd-provider-$*
	cat cmd/herd-provider-example/main.go | sed -e 's/example/$*/g' | gofmt > $@
herd-provider-%: host.go hostset.go go.mod go.sum cmd/herd-provider-%/main.go provider/%/*.go provider/plugin/common/*.go provider/plugin/server/*.go $(protobuf_sources)
	go build $(GOGCFLAGS) -o "$@" github.com/seveas/herd/cmd/$@
provider-plugins-source: $(patsubst %,cmd/herd-provider-%/main.go,$(provider_plugins))
provider-plugins: $(patsubst %,herd-provider-%,$(provider_plugins))

# Generated source files part 1: protobuf
%_grpc.pb.go: %.proto
	protoc --go-grpc_out=. $^

%.pb.go: %.proto
	protoc --go_out=. $^

# Generated source files part 2: antlr
$(antlr_sources): scripting/Herd.g4
	(cd scripting; antlr -Dlanguage=Go -o parser Herd.g4)

# Tests and related targets
lint:
	golangci-lint run ./...

tidy:
	go mod tidy

provider/plugin/testdata/bin/herd-provider-%: host.go hostset.go go.mod go.sum provider/plugin/testdata/provider/%/*.go provider/plugin/testdata/cmd/herd-provider-%/*.go provider/plugin/common/* provider/plugin/server/* $(protobuf_sources)
	go build $(GOGCFLAGS) -o "$@" github.com/seveas/herd/provider/plugin/testdata/cmd/herd-provider-$*

test: test-providers test-go lint tidy test-build provider/plugin/testdata/bin/herd-provider-ci
test-providers: provider/plugin/testdata/bin/herd-provider-ci provider/plugin/testdata/bin/herd-provider-ci_dataloader provider/plugin/testdata/bin/herd-provider-ci_cache
test-go:
	go test ./...
test-build: provider-plugins-source
	GOOS=darwin go build github.com/seveas/herd/cmd/herd
	GOOS=linux go build github.com/seveas/herd/cmd/herd
	GOOS=windows go build github.com/seveas/herd/cmd/herd

ABORT ?= --exit-code-from herd --abort-on-container-exit
test-integration:
	make -C integration/pki
	test -e integration/openssh/user.key || ssh-keygen -t ecdsa -f integration/openssh/user.key -N ""
	docker compose down || true
	docker compose build
	docker compose up $(ABORT)
	docker compose down

# Release mechanism
dist_oses := darwin-amd64 darwin-arm64 dragonfly-amd64 freebsd-amd64 linux-amd64 netbsd-amd64 openbsd-amd64 windows-amd64
VERSION = $(shell go run cmd/version.go)
build-all:
	@echo Building herd
	@$(foreach os,$(dist_oses),echo " - for $(os)" && mkdir -p dist/$(os) && GOOS=$(firstword $(subst -, ,$(os))) GOARCH=$(lastword $(subst -, ,$(os))) go build -tags no_extra -ldflags '-s -w' -o dist/$(os)/herd-$(VERSION)/  github.com/seveas/herd/cmd/herd && tar -C dist/$(os)/ -zcf herd-$(VERSION)-$(os).tar.gz herd-$(VERSION)/;)

clean:
	rm -f herd
	rm -f herd-provider-example
	rm -f $(patsubst %,herd-provider-%,$(provider_plugins))
	rm -f provider/plugin/testdata/bin/herd-provider-ci
	rm -f provider/plugin/testdata/bin/herd-provider-ci_dataloader
	rm -f provider/plugin/testdata/bin/herd-provider-ci_cache
	go mod tidy

fullclean: clean
	rm -rf dist/
	rm -f herd-*.tar.gz

install:
	go install github.com/seveas/herd/cmd/herd

.PHONY: tidy test build-all clean fullclean install test-go test-build test-integration lint external-providers external-providers-source
