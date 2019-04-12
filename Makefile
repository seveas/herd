katyusha: *.go cmd/katyusha/*.go
	go build github.com/seveas/katyusha/cmd/katyusha

fmt:
	gofmt -w *.go cmd/katyusha/*.go
