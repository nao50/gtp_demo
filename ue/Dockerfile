FROM golang
WORKDIR ${GOPATH}/
RUN apt-get update -y && apt-get install -y vim tcpdump iperf3
RUN ip route add 10.0.12.0/24 via 10.0.10.20 dev eth0

RUN go get -u github.com/golang/dep/...

