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
	pb "github.com/casbin/casbin-server/proto"
	"github.com/stretchr/testify/assert"

	"testing"
)

func testEnforce(t *testing.T, e *testEngine, sub string, obj string, act string, res bool) {
	t.Helper()
	reply, err := e.s.Enforce(e.ctx, &pb.EnforceRequest{EnforcerHandler: e.h, Params: []string{sub, obj, act}})
	assert.NoError(t, err)

	if reply.Res != res {
		t.Errorf("%s, %v, %s: %t, supposed to be %t", sub, obj, act, !res, res)
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
