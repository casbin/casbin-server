FROM golang:1.13

RUN apt-get update && \
    apt-get -y install git unzip build-essential autoconf libtool

# Install protobuf from source
RUN git clone https://github.com/google/protobuf.git && \
    cd protobuf && \
    ./autogen.sh && \
    ./configure && \
    make && \
    make install && \
    ldconfig && \
    make clean && \
    cd .. && \
    rm -r protobuf

# Get grpc
RUN go get google.golang.org/grpc

# Install protoc-gen-go
RUN go get github.com/golang/protobuf/protoc-gen-go

# Go environment variable to enable Go modules
ENV GO111MODULE=on

# Copy the source and generate the .proto file
ADD . /go/src/github.com/casbin/casbin-server
WORKDIR $GOPATH/src/github.com/casbin/casbin-server
RUN protoc -I proto --go_out=plugins=grpc:proto proto/casbin.proto

# Download dependencies
RUN go mod download

# Install app
RUN go install .
ENTRYPOINT $GOPATH/bin/casbin-server

EXPOSE 50051
