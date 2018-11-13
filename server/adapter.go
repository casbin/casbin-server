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
	"github.com/casbin/gorm-adapter"
	//_ "github.com/jinzhu/gorm/dialects/mssql"
	//_ "github.com/jinzhu/gorm/dialects/mysql"
	//_ "github.com/jinzhu/gorm/dialects/postgres"
)

var errDriverName = errors.New("currently supported DriverName: file | mysql | postgres | mssql")

func newAdapter(in *pb.NewAdapterRequest) (persist.Adapter, error) {
	var a persist.Adapter
	supportDriverNames := [...]string{"file", "mysql", "postgres", "mssql"}

	switch in.DriverName {
	case "file":
		a = fileadapter.NewAdapter(in.ConnectString)
	default:
		var support = false
		for _, driverName := range supportDriverNames {
			if driverName == in.DriverName {
				support = true
				break
			}
		}
		if support {
			a = gormadapter.NewAdapter(in.DriverName, in.ConnectString, in.DbSpecified)
			break
		}
		return nil, errDriverName
	}

	return a, nil
}
