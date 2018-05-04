FROM golang
WORKDIR ${GOPATH}/src/github.com/tily/sampleapp
RUN apt-get update -y && apt-get install -y vim
RUN go get -u github.com/golang/dep/...