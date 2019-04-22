# Let's not rebuild the parser if we don't have antlr available
ifeq ("", "$(strip $(shell which antlr))")
	antlr_sources :=
else
	antlr_sources := parser/katyusha_base_listener.go parser/katyusha_lexer.go parser/katyusha_listener.go parser/katyusha_parser.go
endif

katyusha: *.go cmd/katyusha/*.go $(antlr_sources)
	go build github.com/seveas/katyusha/cmd/katyusha

$(antlr_sources): Katyusha.g4
	antlr -Dlanguage=Go -o parser Katyusha.g4

fmt:
	gofmt -w . cmd/katyusha

vet:
	go vet github.com/seveas/katyusha github.com/seveas/katyusha/cmd/katyusha

test: fmt vet
