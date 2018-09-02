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

func TestRBACModel(t *testing.T) {
	s := NewServer()
	ctx := context.Background()

	s.NewAdapter(ctx, &pb.NewAdapterRequest{DriverName: "file", ConnectString: "../examples/rbac_policy.csv"})

	modelText, err := ioutil.ReadFile("../examples/rbac_model.conf")
	if err != nil {
		t.Error(err)
	}

	resp, err := s.NewEnforcer(ctx, &pb.NewEnforcerRequest{ModelText: string(modelText), AdapterHandle: 0})
	if err != nil {
		t.Error(err)
	}
	e := resp.Handler

	sub := "alice"
	obj := "data1"
	act := "read"
	res := true

	resp2, err := s.Enforce(ctx, &pb.EnforceRequest{EnforcerHandler: e, Sub: sub, Obj: obj, Act: act})
	if err != nil {
		t.Error(err)
	}
	myRes := resp2.Res

	if myRes != res {
		t.Errorf("%s, %s, %s: %t, supposed to be %t", sub, obj, act, myRes, res)
	}
}
