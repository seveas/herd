FROM golang:1.13

WORKDIR /herd
COPY . .
RUN sed -e /provider/d -i go.mod
RUN make GOFLAGS=-mod=vendor
