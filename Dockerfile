FROM golang:1.10 AS builder

RUN apt-get update && apt-get install protobuf-compiler -y

WORKDIR $GOPATH/src/github.com/r00takaspin/raft

RUN bash -c "go get github.com/golang/protobuf/{proto,protoc-gen-go}"

COPY . ./

RUN make

ENTRYPOINT ["./raft"]

