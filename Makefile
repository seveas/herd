# Let's not rebuild the parser if we don't have antlr available
ifeq ("", "$(strip $(shell which antlr))")
	antlr_sources :=
else
	antlr_sources := parser/katyusha_base_listener.go parser/katyusha_lexer.go parser/katyusha_listener.go parser/katyusha_parser.go
endif

katyusha.bin: *.go katyusha/*.go katyusha/cmd/*.go $(antlr_sources)
	go build -o "$@" github.com/seveas/katyusha/katyusha

$(antlr_sources): Katyusha.g4
	antlr -Dlanguage=Go -o parser Katyusha.g4

fmt:
	go fmt ./...

vet:
	go vet ./...

test: fmt vet
	go test ./...
