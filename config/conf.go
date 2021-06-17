package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const (
	DefaultDB        = "mongodb"
	DefaultDB_USER   = "root"
	DefaultDB_PASS   = "example"
	DefaultURL       = "127.0.0.1"
	RestfulEPDefault = "localhost:8181"
)

type Config struct {
	DBType       string `json:"db_type"`
	DBUser       string `json:"db_user"`
	DBPass       string `json:"db_pass"`
	DBConnection string `json:"db_conn"`
	AddrREST     string `json:"restful_addr"`
}

// ReadConfig reads config file if it exists. Otherwise, returns default configs.
func ReadConfig(filename string) (Config, error) {
	conf := Config{
		DefaultDB,
		DefaultDB_USER,
		DefaultDB_PASS,
		DefaultURL,
		RestfulEPDefault,
	}

	if !strings.HasSuffix(filename, ".json") {
		return conf, fmt.Errorf("cannot open non-json files, filename: %s", filename)
	}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Configuration file not found. Continuing with default values.")
		return conf, err
	}

	err = json.NewDecoder(file).Decode(&conf)
	return conf, err
}
