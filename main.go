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
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "github.com/casbin/casbin-server/proto"
	"github.com/casbin/casbin-server/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// File for liveness probe to watch
const LIVE_FILE = "/tmp/app_server_live"

func main() {
	serveCtx, cancelServeCtx := context.WithCancel(context.Background())

	var port int
	flag.IntVar(&port, "port", 50051, "listening port")
	flag.Parse()

	if port < 1 || port > 65535 {
		panic(fmt.Sprintf("invalid port number: %d", port))
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	svr := server.NewServer()
	pb.RegisterCasbinServer(s, svr)
	// Register reflection service on gRPC server.
	reflection.Register(s)

	err = svr.DefaultEnforcerInit()
	if err != nil {
		log.Fatalf("failed to init default enforcer: %v", err)
	}

	// Create liveness file
	_, err = os.Create(LIVE_FILE)
	if err != nil {
		log.Fatal(err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig
		cleanCtx, cancelCleanCtx := context.WithTimeout(serveCtx, 30*time.Second)

		go func() {
			<-cleanCtx.Done()

			if cleanCtx.Err() == context.DeadlineExceeded {
				log.Fatal("Graceful shutdown timed out... Forcing exit now...")
			}
		}()

		log.Println("Gracefully shutdown server...")

		// Remove liveness file
		err = os.Remove(LIVE_FILE)
		if err != nil {
			log.Fatal(err)
		}

		cancelCleanCtx()
		cancelServeCtx()
	}()

	log.Println("Listening on", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	<-serveCtx.Done()
}
