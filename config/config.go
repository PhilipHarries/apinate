package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
)

type QueryKeyType struct {
	KeyName string
	Default string
}

func (q *QueryKeyType) IsEmpty() bool {
	if q.KeyName == "" {
		return true
	}
	return false
}

type Mapping struct {
	Resource  string
	Command   string
	Params    bool
	QueryKeys []QueryKeyType
	Template  string
}

type Config struct {
	Address     string
	Port        int
	ContentType string
	Logfile     string
	Mappings    []Mapping
}

func LoadConfig(filename string) (Config, error) {
	var config Config
	var err error
	switch {
	case filepath.Ext(filename) == ".toml":
		config, err = LoadConfigTOML(filename)
	case filepath.Ext(filename) == ".json":
		config, err = LoadConfigJSON(filename)
	case filepath.Ext(filename) == ".yaml":
		config, err = LoadConfigYAML(filename)
	}
	if err != nil {
		return config, err
	} else {
		if config.Logfile == "" {
			config.Logfile = "stdout"
		}
		if config.Port == 0 {
			config.Port = 8080
		}
		if config.Address == "" {
			config.Address = "0.0.0.0"
		}
	}
	return config, nil
}

// Open configuration file and decode the TOML
func LoadConfigTOML(filename string) (Config, error) {
	var c Config
	_, err := os.Stat(filename)
	if err != nil {
		return c, err
	}
	if _, err := toml.DecodeFile(filename, &c); err != nil {
		return c, err
	}
	return c, nil
}

// Open configuration file and decode the JSON
func LoadConfigJSON(filename string) (Config, error) {
	var c Config
	f, err := os.Open(filename)
	if err != nil {
		return c, err
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	dec.Decode(&c)
	return c, nil
}

// Open configuration file and decode the YAML
func LoadConfigYAML(filename string) (Config, error) {
	var c Config
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		return c, err
	}
	err = yaml.Unmarshal(f, &c)
	return c, nil
}
