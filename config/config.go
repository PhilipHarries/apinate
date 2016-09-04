package config

import (
	"encoding/json"
	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Mapping struct {
	Resource string
	Command  string
	Params   bool
}

type Config struct {
	ContentType string
	Mappings    []Mapping
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
