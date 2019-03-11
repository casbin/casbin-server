Casbin Server
====

[![Build Status](https://travis-ci.org/casbin/casbin-server.svg?branch=master)](https://travis-ci.org/casbin/casbin-server)
[![Docker](https://img.shields.io/docker/build/casbin/casbin-server.svg)](https://hub.docker.com/r/casbin/casbin-server/builds/)
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
PHP | https://github.com/php-casbin/casbin-client

Contributions for clients in other languages are welcome :)

## Installation

    go get github.com/casbin/casbin-server

## Database Support

Similar to Casbin, Casbin-Server also uses adapters to provide policy storage. However, because Casbin-Server is a service instead of a library, the adapters have to be implemented inside Casbin-Server. As Golang is a static language, each adapter requires to import 3rd-party library for that database. We cannot import all those 3rd-party libraries inside Casbin-Server's code, as it causes dependency overhead.

For now, only [Gorm Adapter](https://github.com/casbin/casbin-server/blob/master/server/adapter.go) is built-in with ``mssql``, ``mysql``, ``postgres`` imports all commented. If you want to use ``Gorm Adapter`` with one of those databases, you should uncomment that import line, or add your own import, or even use another adapter by modifying Casbin-Server's source code.

## Limitation of ABAC

Casbin-Server also supports the ABAC model as the Casbin library does. You may wonder how Casbin-Server passes the Go structs to the server-side via network? Good question. In fact, Casbin-Server's client dumps Go struct into JSON and transmits the JSON string prefixed by ``ABAC::`` to Casbin-Server. Casbin-Server will recognize the prefix and load the JSON string into a pre-defined Go struct with 11 string members, then pass it to Casbin. So there will be two limitations for Casbin-Server's ABAC compared to Casbin's ABAC:

1. The Go struct should be flat, all members should be primitive types, e.g., string, int, boolean. No nested struct, no slice or map.

2. The Go struct is limited to 11 members at most. If you want to have more members, you should modify [Casbin-Server's source code](https://github.com/casbin/casbin-server/blob/5e21d10e863c7d8461f951417eb1c63fa00204fb/server/abac.go#L27-L40) by adding more members and rebuild it.

## Getting Help

- [Casbin](https://github.com/casbin/casbin)

## License

This project is under Apache 2.0 License. See the [LICENSE](LICENSE) file for the full license text.
