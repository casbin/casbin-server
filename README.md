Casbin Server
====

[![Build Status](https://travis-ci.org/casbin/casbin-server.svg?branch=master)](https://travis-ci.org/casbin/casbin-server)
[![Coverage Status](https://coveralls.io/repos/github/casbin/casbin-server/badge.svg?branch=master)](https://coveralls.io/github/casbin/casbin-server?branch=master)
[![Godoc](https://godoc.org/github.com/casbin/casbin-server?status.svg)](https://godoc.org/github.com/casbin/casbin-server)

Casbin Server is the ``Access Control as a Service (ACaaS)`` solution based on [Casbin](https://github.com/casbin/casbin). It provides [gRPC](https://grpc.io/) interface for Casbin authorization.

## What is ``Casbin Server``?

Casbin-Server is just a container of Casbin enforcers and adapters. Casbin-Server is designed to be ``compute-intensive`` (for calculating whether an access should be allowed) instead of a centralized policy storage. Just like how native Casbin library works, each Casbin enforcer in Casbin-Server can use its own adapter, which is linked with external database for policy storage.

Of course, you can setup Casbin-Server together with your policy database in the same machine. But they can be separated. If you want to achieve high availability, you can use a Redis cluster as policy storage, then link Casbin-Server's adapter with it. In this sense, Casbin enforcer can be viewed as stateless component. It just retrieves the policy rules it is interested in (via sharding), does some computation and then returns ``allow`` or ``deny``.

## Architecture

Casbin-Server uses the client-server architecture. Casbin-Server itself is the server (in Golang only for now). The clients for Casbin-Server are listed here:

Language | Client
----|----
Golang | https://github.com/casbin/casbin-go-client

Contributions for clients in other languages are welcome :)

## Prerequisites

### Protobuf (if not installed):

As Casbin-Server uses ``gRPC``, you need to install ProtoBuf's command-line tool: ``protoc`` first to generate the ``.proto`` file. The command is:

```
protoc -I proto --go_out=plugins=grpc:proto proto/casbin.proto
```

Alternatively, you can also [run it from an IDE](https://github.com/casbin/casbin-server/blob/6b46c48c8845dc1b8021f2872be08b8e1a62b092/main.go#L15).

## Installation

    go get github.com/casbin/casbin-server

## Getting Help

- [Casbin](https://github.com/casbin/casbin)

## License

This project is under Apache 2.0 License. See the [LICENSE](LICENSE) file for the full license text.
