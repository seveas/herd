FROM golang:1.22

RUN curl https://apt.releases.hashicorp.com/gpg | gpg --dearmor -o /usr/share/keyrings/hashicorp-archive-keyring.gpg && \
    . /etc/os-release && \
    echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $VERSION_CODENAME main" > /etc/apt/sources.list.d/hashicorp.list && \
    apt-get update && apt-get -y install consul

WORKDIR /herd
COPY . .
RUN touch provider/plugin/common/plugin.pb.go provider/plugin/common/plugin_grpc.pb.go
RUN make
RUN make -C integration/pki install-ca
