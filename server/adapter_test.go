package server

import (
	"os"
	"testing"

	miniredis "github.com/alicebob/miniredis/v2"

	pb "github.com/casbin/casbin-server/proto"
	"github.com/stretchr/testify/assert"
)

func TestGetLocalConfig(t *testing.T) {
	assert.Equal(t, configFileDefaultPath, getLocalConfigPath(), "read from default connection config path if environment variable is not set")

	os.Setenv(configFilePathEnvironmentVariable, "dir/custom_path.json")
	assert.Equal(t, "dir/custom_path.json", getLocalConfigPath())
}

func runFakeRedis(username string, password string) (host string, port string, err error) {
	s, err := miniredis.Run()
	if err != nil {
		return "", "", err
	}
	if username != "" && password != "" {
		s.RequireUserAuth(username, password)
	}
	return s.Host(), s.Port(), err
}

func TestRedisAdapterConfig(t *testing.T) {
	os.Setenv(configFilePathEnvironmentVariable, "../config/connection_config.json")

	host, port, err := runFakeRedis("", "")

	in := &pb.NewAdapterRequest{
		DriverName:    "redis",
		ConnectString: "redis://" + host + ":" + port,
	}

	a, err := newAdapter(in)
	assert.NoError(t, err, "should create redis adapter without error")
	assert.NotNil(t, a, "adapter should not be nil")
}

func TestRedisAdapterConfigWithUsernameAndPassword(t *testing.T) {
	os.Setenv(configFilePathEnvironmentVariable, "../config/connection_config.json")

	username, password := "foo", "bar"
	host, port, err := runFakeRedis(username, password)

	in := &pb.NewAdapterRequest{
		DriverName:    "redis",
		ConnectString: "redis://" + username + ":" + password + "@" + host + ":" + port,
	}

	a, err := newAdapter(in)
	assert.NoError(t, err, "should create redis adapter without error")
	assert.NotNil(t, a, "adapter should not be nil")
}

func TestRedisAdapterConfigWithoutPrefix(t *testing.T) {
	os.Setenv(configFilePathEnvironmentVariable, "../config/connection_config.json")

	host, port, err := runFakeRedis("", "")

	in := &pb.NewAdapterRequest{
		DriverName:    "redis",
		ConnectString: host + ":" + port,
	}

	a, err := newAdapter(in)
	assert.NoError(t, err, "should create redis adapter without error")
	assert.NotNil(t, a, "adapter should not be nil")
}

func TestInvalidRedisAdapterConfig(t *testing.T) {
	os.Setenv(configFilePathEnvironmentVariable, "../config/connection_config.json")

	_, _, err := runFakeRedis("", "")

	in := &pb.NewAdapterRequest{
		DriverName:    "redis",
		ConnectString: "invalid-address",
	}

	a, err := newAdapter(in)
	assert.Error(t, err, "should cause an redis adapter without error")
	assert.ErrorContains(t, err, "dial tcp: lookup invalid-address")
	assert.Nil(t, a, "adapter should be nil")
}

func TestRedisAdapterConfigReturnDefaultFallback(t *testing.T) {
	os.Setenv(configFilePathEnvironmentVariable, "../config/connection_config.json")

	in := &pb.NewAdapterRequest{
		DriverName:    "redis",
		ConnectString: "",
	}

	a, err := newAdapter(in)
	assert.NoError(t, err, "should create file default adapter without error")
	assert.NotNil(t, a, "adapter should not be nil")
}
