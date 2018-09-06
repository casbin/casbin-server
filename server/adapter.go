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
	"errors"
	pb "github.com/casbin/casbin-server/proto"
	"github.com/casbin/casbin/persist"
	"github.com/casbin/casbin/persist/file-adapter"
)

var errDriverName = errors.New("invalid DriverName")

const (
	MYSQL = "mysql"
	FILE  = "file"
)

func newAdapter(in *pb.NewAdapterRequest) (persist.Adapter, error) {
	var a persist.Adapter

	switch in.DriverName {
	case FILE:
		a = fileadapter.NewAdapter(in.ConnectString)
	default:
		return nil, errDriverName
	}

	return a, nil
}
