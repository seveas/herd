FROM golang:1.15

WORKDIR /herd
COPY . .
RUN touch provider/plugin/common/plugin.pb.go provider/plugin/common/plugin_grpc.pb.go
RUN make GOFLAGS=-mod=vendor
RUN make -C integration/pki install-ca
