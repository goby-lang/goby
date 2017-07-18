FROM golang:latest

RUN go get github.com/tools/godep
RUN go get github.com/goby-lang/goby

WORKDIR $GOPATH/src/github.com/goby-lang/goby
