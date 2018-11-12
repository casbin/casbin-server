FROM grpc/go
ADD . /go/src/github.com/casbin/casbin-server
WORKDIR $GOPATH/src/github.com/casbin/casbin-server
RUN protoc -I proto --go_out=plugins=grpc:proto proto/casbin.proto

# Download and install the latest release of dep
ADD https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep

# Install dependencies
RUN dep init
RUN dep ensure --vendor-only

# Install app
RUN go install .
ENTRYPOINT $GOPATH/bin/casbin-server

EXPOSE 50051
