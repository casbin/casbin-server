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
	"os"
	"testing"

	pb "github.com/casbin/casbin-server/proto"
	"github.com/stretchr/testify/assert"
)

func testEnforce(t *testing.T, e *testEngine, sub string, obj string, act string, res bool) {
	t.Helper()
	reply, err := e.s.Enforce(e.ctx, &pb.EnforceRequest{EnforcerHandler: e.h, Params: []string{sub, obj, act}})
	assert.NoError(t, err)

	if reply.Res != res {
		t.Errorf("%s, %v, %s: %t, supposed to be %t", sub, obj, act, !res, res)
	} else {
		t.Logf("Enforce for %s, %s, %s : %v", sub, obj, act, reply.Res)
	}
}

func testEnforceWithoutUsers(t *testing.T, e *testEngine, obj string, act string, res bool) {
	t.Helper()
	reply, err := e.s.Enforce(e.ctx, &pb.EnforceRequest{EnforcerHandler: e.h, Params: []string{obj, act}})
	assert.NoError(t, err)

	if reply.Res != res {
		t.Errorf("%s, %s: %t, supposed to be %t", obj, act, !res, res)
	}
}

func TestRBACModel(t *testing.T) {
	s := NewServer()
	ctx := context.Background()

	_, err := s.NewAdapter(ctx, &pb.NewAdapterRequest{DriverName: "file", ConnectString: "../examples/rbac_policy.csv"})
	if err != nil {
		t.Error(err)
	}

	modelText, err := os.ReadFile("../examples/rbac_model.conf")
	if err != nil {
		t.Error(err)
	}

	resp, err := s.NewEnforcer(ctx, &pb.NewEnforcerRequest{ModelText: string(modelText), AdapterHandle: 0, EnableAcceptJsonRequest: false})
	if err != nil {
		t.Error(err)
	}
	e := resp.Handler

	sub := "alice"
	obj := "data1"
	act := "read"
	res := true

	resp2, err := s.Enforce(ctx, &pb.EnforceRequest{EnforcerHandler: e, Params: []string{sub, obj, act}})
	if err != nil {
		t.Error(err)
	}
	myRes := resp2.Res

	if myRes != res {
		t.Errorf("%s, %s, %s: %t, supposed to be %t", sub, obj, act, myRes, res)
	}
}

func TestABACModel(t *testing.T) {
	s := NewServer()
	ctx := context.Background()

	modelText, err := os.ReadFile("../examples/abac_model.conf")
	if err != nil {
		t.Error(err)
	}

	resp, err := s.NewEnforcer(ctx, &pb.NewEnforcerRequest{ModelText: string(modelText), AdapterHandle: -1, EnableAcceptJsonRequest: false})
	if err != nil {
		t.Error(err)
	}
	type ABACModel struct {
		Name  string
		Owner string
	}
	e := resp.Handler

	data1, _ := MakeABAC(ABACModel{Name: "data1", Owner: "alice"})
	data2, _ := MakeABAC(ABACModel{Name: "data2", Owner: "bob"})

	testModel(t, s, e, "alice", data1, "read", true)
	testModel(t, s, e, "alice", data1, "write", true)
	testModel(t, s, e, "alice", data2, "read", false)
	testModel(t, s, e, "alice", data2, "write", false)
	testModel(t, s, e, "bob", data1, "read", false)
	testModel(t, s, e, "bob", data1, "write", false)
	testModel(t, s, e, "bob", data2, "read", true)
	testModel(t, s, e, "bob", data2, "write", true)

}

func testModel(t *testing.T, s *Server, enforcerHandler int32, sub string, obj string, act string, res bool) {
	t.Helper()

	reply, err := s.Enforce(context.TODO(), &pb.EnforceRequest{EnforcerHandler: enforcerHandler, Params: []string{sub, obj, act}})
	assert.NoError(t, err)

	if reply.Res != res {
		t.Errorf("%s, %v, %s: %t, supposed to be %t", sub, obj, act, !res, res)
	}
}
