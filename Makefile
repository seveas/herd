herd: *.go cmd/herd/*.go
	go build github.com/seveas/herd/cmd/herd

fmt:
	gofmt -w *.go cmd/herd/*.go
