package server

import (
	"os"
	"testing"

	pb "github.com/casbin/casbin-server/proto"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	"github.com/stretchr/testify/assert"
)

func TestGetLocalConfig(t *testing.T) {
	assert.Equal(t, configFileDefaultPath, getLocalConfigPath(), "read from default connection config path if environment variable is not set")

	os.Setenv(configFilePathEnvironmentVariable, "dir/custom_path.json")
	assert.Equal(t, "dir/custom_path.json", getLocalConfigPath())
}
func TestEmptyAdapter(t *testing.T) {

	// Setup
	in1 := &pb.NewAdapterRequest{DriverName: "", ConnectString: ""}
	a, err2 := newAdapter(in1)
	assert.Error(t, err2)
	assert.Nil(t, a)
}
func TestNewAdapter(t *testing.T) {

	// Setup
	in1 := &pb.NewAdapterRequest{DriverName: "file", ConnectString: "path/to/file"}
	a1 := fileadapter.NewAdapter(in1.ConnectString)
	a, err := newAdapter(in1)
	assert.NoError(t, err)
	assert.Equal(t, a1, a)
	a2, err2 := newAdapter(in1)
	assert.Error(t, err2)
	assert.Equal(t, a, a2)

	// Test case 2: invalid driver name
	in2 := &pb.NewAdapterRequest{DriverName: "invalid", ConnectString: "invalid"}
	_, err = newAdapter(in2)
	assert.Error(t, err)
}
