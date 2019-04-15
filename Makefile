katyusha: *.go cmd/katyusha/*.go
	go build github.com/seveas/katyusha/cmd/katyusha

fmt:
	gofmt -w . cmd/katyusha

vet:
	go vet github.com/seveas/katyusha github.com/seveas/katyusha/cmd/katyusha

test: fmt vet
