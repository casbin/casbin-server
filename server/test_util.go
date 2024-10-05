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

package server

import (
	"context"
	"io/ioutil"
	"testing"

	pb "github.com/casbin/casbin-server/proto"
)

type testEngine struct {
	s   *Server
	ctx context.Context
	h   int32
}

func newTestEngine(t *testing.T, from, connectStr string, modelLoc string) *testEngine {
	s := NewServer()
	ctx := context.Background()

	_, err := s.NewAdapter(ctx, &pb.NewAdapterRequest{DriverName: from, ConnectString: connectStr})
	if err != nil {
		t.Fatal(err)
	}

	modelText, err := ioutil.ReadFile(modelLoc)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := s.NewEnforcer(ctx, &pb.NewEnforcerRequest{ModelText: string(modelText), AdapterHandle: 0, EnableAcceptJsonRequest: true})
	if err != nil {
		t.Fatal(err)
	}

	return &testEngine{s: s, ctx: ctx, h: resp.Handler}
}
