FROM golang:latest
WORKDIR ${GOPATH}/src/
RUN apt-get update -y && apt-get install -y vim tcpdump
RUN go get -u github.com/golang/dep/...