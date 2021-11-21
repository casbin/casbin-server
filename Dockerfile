FROM golang:1.16

RUN apt-get update && \
    apt-get -y install git unzip build-essential autoconf libtool

# Install protobuf from source
RUN git clone --depth=1 https://github.com/protocolbuffers/protobuf.git && \
    cd protobuf && \
    ./autogen.sh && \
    ./configure && \
    make && \
    make install && \
    ldconfig && \
    make clean && \
    cd .. && \
    rm -r protobuf

# Go environment variable to enable Go modules
ENV GO111MODULE=on

# Get grpc
RUN go get google.golang.org/grpc

# Install protoc-gen-go
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Copy the source and generate the .proto file
ADD . /go/src/github.com/casbin/casbin-server
WORKDIR $GOPATH/src/github.com/casbin/casbin-server
RUN protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=require_unimplemented_servers=false \
    --go-grpc_opt=paths=source_relative proto/casbin.proto

# Download dependencies
RUN go mod download

# Install app
RUN go install .
ENTRYPOINT $GOPATH/bin/casbin-server

EXPOSE 50051
