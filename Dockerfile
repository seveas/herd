FROM golang:1.13

WORKDIR /katyusha
COPY . .
RUN sed -e /provider/d -i go.mod
RUN make
