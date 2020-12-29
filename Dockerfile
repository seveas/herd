FROM golang:1.15

WORKDIR /herd
COPY . .
RUN make GOFLAGS=-mod=vendor
