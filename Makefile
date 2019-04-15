herd: *.go cmd/herd/*.go
	go build github.com/seveas/herd/cmd/herd

fmt:
	gofmt -w . cmd/herd

vet:
	go vet github.com/seveas/herd github.com/seveas/herd/cmd/herd

test: fmt vet
