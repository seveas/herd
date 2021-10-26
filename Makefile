# Let's not rebuild the parser if we don't have antlr available
ifeq ("", "$(strip $(shell which antlr))")
	antlr_sources :=
else
	antlr_sources := scripting/parser/katyusha_base_listener.go scripting/parser/katyusha_lexer.go scripting/parser/katyusha_listener.go scripting/parser/katyusha_parser.go
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

katyusha: go.mod go.sum *.go cmd/katyusha/*.go sshagent/*.go scripting/*.go provider/*/*.go provider/plugin/common/*.go $(protobuf_sources) $(antlr_sources)
	go build -o "$@" github.com/seveas/katyusha/cmd/katyusha

%_grpc.pb.go: %.proto
	protoc --go-grpc_out=. $^

%.pb.go: %.proto
	protoc --go_out=. $^

katyusha-provider-%: go.mod go.sum cmd/katyusha-provider-%/*.go provider/%/*.go provider/plugin/common/* provider/plugin/server/* $(protobuf_sources)
	go build -o "$@" github.com/seveas/katyusha/cmd/$@

$(antlr_sources): scripting/Katyusha.g4
	(cd scripting; antlr -Dlanguage=Go -o parser Katyusha.g4)

fmt:
	go fmt ./...

vet:
	go vet ./...

tidy:
	go mod tidy

provider/plugin/testdata/bin/katyusha-provider-ci: go.mod go.sum provider/plugin/testdata/provider/ci/*.go provider/plugin/testdata/cmd/katyusha-provider-ci/*.go provider/plugin/common/* provider/plugin/server/* $(protobuf_sources)
	go build -o "$@" github.com/seveas/katyusha/provider/plugin/testdata/cmd/katyusha-provider-ci

test: fmt vet tidy provider/plugin/testdata/bin/katyusha-provider-ci
	go test ./...
	GOOS=windows go build github.com/seveas/katyusha/cmd/katyusha

ABORT ?= --exit-code-from katyusha --abort-on-container-exit
test-integration:
	go mod vendor
	make -C integration/pki
	test -e integration/openssh/user.key || ssh-keygen -t ecdsa -f integration/openssh/user.key -N ""
	docker-compose down || true
	docker-compose build
	docker-compose up $(ABORT)
	docker-compose down

dist_oses := darwin dragonfly freebsd linux netbsd openbsd windows
ssh_agent_oses := darwin dragonfly freebsd linux netbsd openbsd
build_all:
	@echo Building katyusha
	@$(foreach os,$(dist_oses),echo " - for $(os)" && mkdir -p dist/$(os)-amd64 && GOOS=$(os) GOARCH=amd64 go build -tags no_extra -ldflags '-s -w' -o dist/$(os)-amd64/ github.com/seveas/katyusha/cmd/katyusha;)

clean:
	rm -f katyusha
	rm -f katyusha-provider-example
	go mod tidy

fullclean: clean
	rm -rf dist/
	rm -f $(antlr_sources)

install:
	go install github.com/seveas/katyusha/cmd/katyusha

.PHONY: fmt vet tidy test build_all clean fullclean install
