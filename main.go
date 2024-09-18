// Copyright 2018 The casbin Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:generate protoc -I proto --go_out=plugins=grpc:proto proto/casbin.proto

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"sync"

	ginHandler "github.com/casbin/casbin-server/handler/gin"
	pb "github.com/casbin/casbin-server/proto"
	"github.com/casbin/casbin-server/router"
	ginRouter "github.com/casbin/casbin-server/router/gin"
	"github.com/casbin/casbin-server/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		var port int
		flag.IntVar(&port, "port", 50051, "gRPC listening port")
		flag.Parse()

		if port < 1 || port > 65535 {
			panic(fmt.Sprintf("invalid gRPC port number: %d", port))
		}

		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		pb.RegisterCasbinServer(s, server.NewServer())
		// Register reflection service on gRPC server.
		reflection.Register(s)
		log.Println("gRPC listening on", port)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to gRPC serve: %v", err)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		var port int
		flag.IntVar(&port, "http-port", 8585, "http listening port")
		flag.Parse()
		if port < 1 || port > 65535 {
			panic(fmt.Sprintf("invalid http port number: %d", port))
		}
		var r router.Router
		r = ginRouter.New() // or echoRouter.New()
		// Define handlers
		h := ginHandler.NewHttpHandler()
		r.POST("/authorize", h.Enforce)

		// Start the server
		httpAddr := fmt.Sprintf(":%d", port)
		log.Println("http listening on", httpAddr)
		if err := r.Serve(httpAddr); err != nil {
			log.Fatalf("failed to http serve: %v", err)
		}

	}()
	wg.Wait()
}
