FROM golang:1.17

RUN apt-get update && apt-get -y install consul

WORKDIR /herd
COPY . .
RUN touch provider/plugin/common/plugin.pb.go provider/plugin/common/plugin_grpc.pb.go
RUN make
RUN make -C integration/pki install-ca
