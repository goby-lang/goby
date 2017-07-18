FROM golang:latest

ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH

RUN go get github.com/tools/godep
RUN go get github.com/lib/pq
RUN go build -buildmode=plugin -linkshared std github.com/lib/pq database/sql

RUN mkdir -p $GOPATH/src/github.com/goby-lang/goby

WORKDIR $GOPATH/src/github.com/goby-lang/goby

ADD . ./
RUN godep restore

RUN ls $GOPATH/src/github.com/lib/pq
#RUN go run goby.go ./samples/import.gb
