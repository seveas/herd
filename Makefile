# Let's not rebuild the parser if we don't have antlr available
ifeq ("", "$(strip $(shell which antlr))")
	antlr_sources :=
else
	antlr_sources := parser/herd_base_listener.go parser/herd_lexer.go parser/herd_listener.go parser/herd_parser.go
endif

herd.bin: go.mod *.go herd/*.go herd/cmd/*.go $(antlr_sources)
	go build -o "$@" github.com/seveas/herd/herd

$(antlr_sources): Herd.g4
	antlr -Dlanguage=Go -o parser Herd.g4

fmt:
	go fmt ./...

vet:
	go vet ./...

test: fmt vet
	go test ./...

dist_oses := darwin dragonfly freebsd linux netbsd openbsd windows
build_all:
	$(foreach os,$(dist_oses),echo "Building for $(os)" && mkdir -p dist/$(os)-amd64 && GOOS=$(os) GOARCH=amd64 go build -ldflags '-s -w' -o dist/$(os)-amd64/ github.com/seveas/herd/herd;)
