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
	gofmt -w . katyusha katyusha/cmd

vet:
	go vet github.com/seveas/katyusha github.com/seveas/katyusha/katyusha github.com/seveas/katyusha/katyusha/cmd

test: fmt vet
