FROM golang:latest

ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH

RUN apt-get update && apt-get install -y zsh

RUN go get github.com/tools/godep
RUN go build -buildmode=plugin std

RUN mkdir -p $GOPATH/src/github.com/goby-lang/goby

WORKDIR $GOPATH/src/github.com/goby-lang/goby

RUN mkdir plugins
RUN mkdir -p Godeps/_workspace
ADD Godeps/Godeps.json ./Godeps

RUN godep restore

ADD . ./

#RUN go build -buildmode=plugin -o ./plugin.so ./plugin/plugin.go

#RUN ls $GOPATH/src/github.com/lib/pq
RUN go run goby.go ./samples/import.gb
