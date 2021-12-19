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

herd: go.mod go.sum *.go cmd/herd/*.go ssh/*.go scripting/*.go provider/*/*.go provider/plugin/common/*.go $(protobuf_sources) $(antlr_sources)
	go build -o "$@" github.com/seveas/herd/cmd/herd

%_grpc.pb.go: %.proto
	protoc --go-grpc_out=. $^

%.pb.go: %.proto
	protoc --go_out=. $^

herd-provider-%: go.mod go.sum cmd/herd-provider-%/*.go provider/%/*.go provider/plugin/common/* provider/plugin/server/* $(protobuf_sources)
	go build -o "$@" github.com/seveas/herd/cmd/$@

$(antlr_sources): scripting/Herd.g4
	(cd scripting; antlr -Dlanguage=Go -o parser Herd.g4)

fmt:
	go fmt ./...

vet:
	go vet ./...

tidy:
	go mod tidy

provider/plugin/testdata/bin/herd-provider-ci: go.mod go.sum provider/plugin/testdata/provider/ci/*.go provider/plugin/testdata/cmd/herd-provider-ci/*.go provider/plugin/common/* provider/plugin/server/* $(protobuf_sources)
	go build -o "$@" github.com/seveas/herd/provider/plugin/testdata/cmd/herd-provider-ci

test: fmt vet tidy provider/plugin/testdata/bin/herd-provider-ci
	go test ./...
	GOOS=windows go build github.com/seveas/herd/cmd/herd

ABORT ?= --exit-code-from herd --abort-on-container-exit
test-integration:
	go mod vendor
	make -C integration/pki
	test -e integration/openssh/user.key || ssh-keygen -t ecdsa -f integration/openssh/user.key -N ""
	docker-compose down || true
	docker-compose build
	docker-compose up $(ABORT)
	docker-compose down

dist_oses := darwin dragonfly freebsd linux netbsd openbsd windows
VERSION = $(shell go run cmd/version.go)
build_all:
	@echo Building herd
	@$(foreach os,$(dist_oses),echo " - for $(os)" && mkdir -p dist/$(os)-amd64 && GOOS=$(os) GOARCH=amd64 go build -tags no_extra -ldflags '-s -w' -o dist/$(os)-amd64/herd-$(VERSION)/  github.com/seveas/herd/cmd/herd && tar -C dist/$(os)-amd64/ -zcvf herd-$(VERSION)-$(os)-amd64.tar.gz herd-$(VERSION)/;)

clean:
	rm -f herd
	rm -f herd-provider-example
	go mod tidy

fullclean: clean
	rm -rf dist/
	rm -f $(antlr_sources)

install:
	go install github.com/seveas/herd/cmd/herd

.PHONY: fmt vet tidy test build_all clean fullclean install
