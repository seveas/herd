# Let's not rebuild the parser if we don't have antlr available
ifeq ("", "$(strip $(shell which antlr))")
	antlr_sources :=
else
	antlr_sources := parser/herd_base_listener.go parser/herd_lexer.go parser/herd_listener.go parser/herd_parser.go
endif

herd: *.go cmd/herd/*.go $(antlr_sources)
	go build github.com/seveas/herd/cmd/herd

$(antlr_sources): Herd.g4
	antlr -Dlanguage=Go -o parser Herd.g4

fmt:
	gofmt -w . cmd/herd

vet:
	go vet github.com/seveas/herd github.com/seveas/herd/cmd/herd

test: fmt vet
