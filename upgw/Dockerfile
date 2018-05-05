FROM golang:latest
WORKDIR ${GOPATH}/
RUN apt-get update -y && apt-get install -y vim tcpdump iperf3
RUN go get -u github.com/golang/dep/...

ADD . ${GOPATH}
