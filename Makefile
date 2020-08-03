# Let's not rebuild the parser if we don't have antlr available
ifeq ("", "$(strip $(shell which antlr))")
	antlr_sources :=
else
	antlr_sources := scripting/parser/herd_base_listener.go scripting/parser/herd_lexer.go scripting/parser/herd_listener.go scripting/parser/herd_parser.go
endif

ifneq ($(origin HERD_TAGS), undefined)
	TAGS := -tags $(HERD_TAGS)
endif

ifeq ($(origin HERD_EXTRA_PROVIDERS), undefined)
	tpp_sources :=
else
	tpp_sources := cmd/herd/extra_providers.go
endif

herd: go.mod *.go cmd/herd/*.go scripting/*.go $(antlr_sources) $(tpp_sources)
	go build $(TAGS) -o "$@" github.com/seveas/herd/cmd/herd

ssh-agent-proxy: go.mod cmd/ssh-agent-proxy/*.go
	go build -o "$@" github.com/seveas/herd/cmd/ssh-agent-proxy

$(antlr_sources): scripting/Herd.g4
	(cd scripting; antlr -Dlanguage=Go -o parser Herd.g4)

$(tpp_sources):
	@echo "Enabling third party providers: $(HERD_EXTRA_PROVIDERS)"
	@echo "package main" > $@
	@echo "import (" >> $@
	@for provider in $(HERD_EXTRA_PROVIDERS); do echo "	_ \"$$provider\"" >>$@; done
	@echo ")" >> $@

fmt:
	go fmt ./...

vet:
	go vet ./...

tidy:
	go mod tidy

test: fmt vet tidy
	go test ./...
	docker-compose down || true
	docker-compose build
	make -C testdata/pki
	docker-compose up --exit-code-from herd --abort-on-container-exit
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
ssh_agent_oses := darwin dragonfly freebsd linux netbsd openbsd
build_all:
	@echo Building herd
	@$(foreach os,$(dist_oses),echo " - for $(os)" && mkdir -p dist/$(os)-amd64 && GOOS=$(os) GOARCH=amd64 go build -tags no_extra -ldflags '-s -w' -o dist/$(os)-amd64/ github.com/seveas/herd/cmd/herd;)
	@echo Building ssh-agent-proxy
	@$(foreach os,$(ssh_agent_oses),echo " - for $(os)" && GOOS=$(os) GOARCH=amd64 go build -ldflags '-s -w' -o dist/$(os)-amd64/ github.com/seveas/herd/cmd/ssh-agent-proxy;)

clean:
	rm -f herd
	rm -f ssh-agent-proxy
	rm -f cmd/herd/extra_providers.go

fullclean: clean
	rm -rf dist/
	rm -f $(antlr_sources)
