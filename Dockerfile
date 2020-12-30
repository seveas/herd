FROM golang:1.15

WORKDIR /katyusha
COPY . .
RUN make GOFLAGS=-mod=vendor
RUN make -C integration/pki install-ca
