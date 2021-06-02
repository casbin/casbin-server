FROM golang:1.16 as builder

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

# Build app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -o casbin-server


FROM alpine:3.13
WORKDIR /
COPY --from=builder /go/src/github.com/casbin/casbin-server/casbin-server ./casbin-server
USER 65532:65532
EXPOSE 50051
ENTRYPOINT ./casbin-server