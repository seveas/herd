FROM golang:1.13

WORKDIR /katyusha
COPY . .
RUN make
