FROM golang:latest

ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH

RUN go get github.com/tools/godep

RUN mkdir -p $GOPATH/src/github.com/goby-lang/goby
ENV GOBY_ROOT=$GOPATH/src/github.com/goby-lang/goby

WORKDIR $GOPATH/src/github.com/goby-lang/goby

RUN mkdir Godeps/
ADD Godeps/Godeps.json ./Godeps

RUN godep restore

ADD . ./

RUN go install .
