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

package main

import (
	"log"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "github.com/casbin/casbin-server/proto"
)

const (
	address     = "localhost:50051"
	defaultName = "Casbin"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewCasbinClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.NewEnforcer(ctx, &pb.NewEnforcerRequest{"", 0})
	if err != nil {
		log.Fatalf("NewEnforcer() error: %v", err)
	}
	log.Printf("Enforer handle: %d", r.Handler)
}
