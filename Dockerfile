FROM golang:1.10 AS builder

RUN apt-get update && apt-get install protobuf-compiler -y

WORKDIR $GOPATH/src/github.com/r00takaspin/raft

RUN bash -c "go get github.com/golang/protobuf/{proto,protoc-gen-go}"

COPY . ./

ADD https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep

RUN dep ensure --vendor-only

RUN make

ENTRYPOINT ["./raft"]

