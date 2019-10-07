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
	"encoding/json"
	"errors"
	"fmt"
	"os"

	pb "github.com/casbin/casbin-server/proto"
	"github.com/casbin/casbin/v2/persist"
	"github.com/casbin/casbin/v2/persist/file-adapter"
	"github.com/casbin/gorm-adapter/v2"
	//_ "github.com/jinzhu/gorm/dialects/mssql"
	//_ "github.com/jinzhu/gorm/dialects/mysql"
	//_ "github.com/jinzhu/gorm/dialects/postgres"
)

var errDriverName = errors.New("currently supported DriverName: file | mysql | postgres | mssql")

func newAdapter(in *pb.NewAdapterRequest) (persist.Adapter, error) {
	var a persist.Adapter
	in = checkLocalConfig(in)
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
		if !support {
			return nil, errDriverName
		}

		var err error
		a, err = gormadapter.NewAdapter(in.DriverName, in.ConnectString, in.DbSpecified)
		if err != nil {
			return nil, err
		}
	}

	return a, nil
}

func checkLocalConfig(in *pb.NewAdapterRequest) *pb.NewAdapterRequest {
	cfg := LoadConfiguration("config/connection_config.json")
	if in.ConnectString == "" || in.DriverName == "" {
		in.DriverName = cfg.Driver
		in.ConnectString = cfg.Connection
	}
	return in
}

func LoadConfiguration(file string) Config {
	//Loads a default config from adapter_config in case a custom adapter isn't provided by the client.
	//DriverName and ConnectionString can be configured in the file. Defaults to 'file' mode.

	configFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err.Error())
	}
	decoder := json.NewDecoder(configFile)
	config := Config{}
	decoder.Decode(&config)
	return config
}

type Config struct {
	Driver string
	Connection string
	Enforcer string
}
