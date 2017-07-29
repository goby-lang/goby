FROM golang:latest

ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH

RUN apt-get update && apt-get install -y zsh

RUN go get github.com/tools/godep

RUN mkdir -p $GOPATH/src/github.com/goby-lang/goby
ENV GOBY_ROOT=$GOPATH/src/github.com/goby-lang/goby

WORKDIR $GOPATH/src/github.com/goby-lang/goby

RUN mkdir Godeps/
ADD Godeps/Godeps.json ./Godeps

RUN godep restore

ADD . ./

# Run test when building image is not a good practice, but it's more convenient for development

#RUN ./test.sh
RUN TEST_PLUGIN=true go test ./vm --run .?Plugin.? -v
RUN TEST_PLUGIN=true go test ./vm --run .?Struct.? -v
