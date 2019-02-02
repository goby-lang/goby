FROM golang:1.11

ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

RUN mkdir -p $GOPATH/src/github.com/gooby-lang/gooby
ENV GOBY_ROOT=$GOPATH/src/github.com/gooby-lang/gooby

WORKDIR $GOPATH/src/github.com/gooby-lang/gooby

ADD . ./

RUN dep ensure

RUN go install .
