package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

type Config struct {
	DefaultUrl string `yaml:"default_url"`
	Triggers   map[string][]string
}

type Trigger struct {
	Path    string
	Actions map[string][]string
}

func (c Config) Load() Config {
	return loadYaml("./dropboy.yml")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func loadYaml(configFile string) Config {
	var config Config

	filename, _ := filepath.Abs(configFile)
	source, err := ioutil.ReadFile(filename)
	check(err)

	err = yaml.Unmarshal(source, &config)
	check(err)

	return config
}
