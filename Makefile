# Let's not rebuild the parser if we don't have antlr available
ifeq ("", "$(strip $(shell which antlr))")
	antlr_sources :=
else
	antlr_sources := scripting/parser/katyusha_base_listener.go scripting/parser/katyusha_lexer.go scripting/parser/katyusha_listener.go scripting/parser/katyusha_parser.go
endif

katyusha: go.mod *.go cmd/katyusha/*.go scripting/*.go $(antlr_sources)
	go build -o "$@" github.com/seveas/katyusha/cmd/katyusha

$(antlr_sources): scripting/Katyusha.g4
	(cd scripting; antlr -Dlanguage=Go -o parser Katyusha.g4)

fmt:
	go fmt ./...

vet:
	go vet ./...

test: fmt vet
	go test ./...
	docker-compose down || true
	docker-compose build
	make -C testdata/pki
	docker-compose up --exit-code-from katyusha --abort-on-container-exit
	docker-compose down

test-integration:
	(cd /etc/ssl/certs && ln -sf ca.crt $$(openssl x509 -in ca.crt -hash -noout).crt)
	cd integration; for f in t*.sh; do \
		if [ -f "$$f" ]; then \
			echo "$$f"; \
			if  ! sh "$$f"; then \
				sh "$$f" --verbose; \
				ec=1; \
			fi; \
		fi; \
	done; exit $$ec

dist_oses := darwin dragonfly freebsd linux netbsd openbsd windows
build_all:
	$(foreach os,$(dist_oses),echo "Building for $(os)" && mkdir -p dist/$(os)-amd64 && GOOS=$(os) GOARCH=amd64 go build -ldflags '-s -w' -o dist/$(os)-amd64/ github.com/seveas/katyusha/cmd/katyusha;)
