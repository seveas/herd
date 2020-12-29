FROM golang:1.15

WORKDIR /katyusha
COPY . .
RUN make GOFLAGS=-mod=vendor
