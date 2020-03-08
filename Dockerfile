FROM golang:1.14

ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH

ENV GO111MODULE=on

RUN mkdir -p $GOPATH/src/github.com/goby-lang/goby
ENV GOBY_ROOT=$GOPATH/src/github.com/goby-lang/goby

WORKDIR $GOPATH/src/github.com/goby-lang/goby

ADD . ./

RUN go install .
