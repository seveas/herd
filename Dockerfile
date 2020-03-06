FROM golang:1.13

WORKDIR /herd
COPY . .
RUN make
