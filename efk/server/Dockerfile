FROM golang:latest AS builder
WORKDIR /server
ENV GO111MODULE=on

COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build ./server/*.go
RUN ls
CMD [ "./main" ]