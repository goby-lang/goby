FROM golang:1.10

ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

RUN mkdir -p $GOPATH/src/github.com/goby-lang/goby
ENV GOBY_ROOT=$GOPATH/src/github.com/goby-lang/goby

WORKDIR $GOPATH/src/github.com/goby-lang/goby

ADD . ./

RUN dep ensure

RUN go install .
