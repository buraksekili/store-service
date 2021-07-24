package persistence

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/buraksekili/store-service/config"
	e "github.com/pkg/errors"
)

var (
	errMissingEnvironmentVariable = errors.New("failed to parse environment variable")
	errInvalidEnvironmentVariable = errors.New("invalid environment variable")
)

const (
	defDBNameEnv = "S_MONGODB_DB"
	defDBHostEnv = "S_MONGODB_HOST"
	defDBPortEnv = "S_MONGODB_PORT"
)

// MongoParser represents a environment variable parser for MongoDB related
// environment variables. All environment variables are read from ./docker/.env
// file. If you would like to configure their default value, check .env.
type MongoParser struct {
	DBName string
	DBHost string
	DBPort string
}

// NewMongoConfigParser returns a ConfigParser that parses environment variables
// regarding MongoDB layer.
func NewMongoConfigParser() config.ConfigParser {
	return &MongoParser{}
}

// Parse parses environment variables regarding MongoDB layer.
// It returns a configuration struct that consists of required
// configuration variables to connect a database. Returns non-nil
// error, if any of the environment variables are equal to the
// empty string.
func (mc *MongoParser) Parse() error {
	if v := os.Getenv(defDBNameEnv); v != "" {
		mc.DBName = v
	} else {
		return e.Wrap(fmt.Errorf("environment variable: %s", defDBNameEnv), errMissingEnvironmentVariable.Error())
	}

	if v := os.Getenv(defDBHostEnv); v != "" {
		mc.DBHost = v
	} else {
		return e.Wrap(fmt.Errorf("environment variable: %s", defDBHostEnv), errMissingEnvironmentVariable.Error())
	}

	if v := os.Getenv(defDBPortEnv); v != "" {
		mc.DBPort = v
	} else {
		return e.Wrap(fmt.Errorf("environment variable: %s", defDBPortEnv), errMissingEnvironmentVariable.Error())
	}

	return nil
}

// Address returns an URL address to be able to connect
// the persistence layer.
func (mc *MongoParser) Address() (string, error) {
	if !valid(mc.DBPort, mc.DBHost, mc.DBName) {
		return "", errInvalidEnvironmentVariable
	}
	return fmt.Sprintf("mongodb://%s:%s", mc.DBHost, mc.DBPort), nil
}

func valid(strs ...string) bool {
	for _, s := range strs {
		if len(strings.TrimSpace(s)) == 0 {
			return false
		}
	}
	return true
}
