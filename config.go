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

func (c Config) Load(args ...interface{}) Config {
  configFile := "./dropboy.yml"

  // an example of how to do default parameters in go
  // http://joneisen.tumblr.com/post/53695478114/golang-and-default-values
  for _, arg := range args {
    switch t := arg.(type) {
    case string:
      configFile = t
    default:
      panic("Unknown argument")
    }
  }

	return loadYaml(configFile)
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
