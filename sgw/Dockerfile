FROM golang
WORKDIR ${GOPATH}/
RUN apt-get update -y && apt-get install -y vim tcpdump iperf3
RUN go get -u github.com/golang/dep/...

ADD . ${GOPATH}

CMD ["go", "run", "sgw.go"]
