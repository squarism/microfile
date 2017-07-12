package config

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/hashicorp/hcl"
	"github.com/mitchellh/go-homedir" // avoids cgo cross compile issues
	"github.com/spf13/viper"          // leaning heavily on viper for configuration
)

type Config struct {
	DefaultURL string  `hcl:"default_url"`
	Watches    []Watch `hcl:"watch"`
}

type Watch struct {
	Path    string   `hcl:",key"`
	Actions []Action `hcl:"action"`
}

type Action struct {
	Type    string            `hcl:",key"`
	Options map[string]string `hcl:"options" hcle:"omitempty"`
}

// read in config defaults and populate our config struct
func populateConfig(c *Config) {
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Can't read config file: %s \n", err))
	}

	// because viper doesn't pass HCL flags down for auto-key niceness
	// lets just bypass viper a bit
	viperConfigFileContents, err := ioutil.ReadFile(viper.GetViper().ConfigFileUsed())
	if err != nil {
		panic(fmt.Errorf("Can't read file, %v", err))
	}

	// read everything into the config reference (c), this is like &config from the HCL examples
	err = hcl.Decode(c, string(viperConfigFileContents))
	if err != nil {
		panic(fmt.Errorf("Problem with the HCL config file, %v", err))
	}
}

func (c *Config) homeConfigDirectory() string {
	home, err := homedir.Dir()
	if err != nil {
		log.Panic(err)
	}

	// This will allow users to put the config in ~/.dropboy/dropboy.yml
	homeConfigDirectory := fmt.Sprintf("%s/.dropboy/", home)
	return homeConfigDirectory
}

// set config options and default values in the config
func setConfigurationDefaults() {
	// viper.SetDefault("DefaultURL", "http://localhost:3000")
	viper.SetConfigName("dropboy")
	viper.SetConfigType("hcl")
}

// Configure uses viper to manage configuration
func (c *Config) Configure(configPaths ...string) {
	setConfigurationDefaults()

	// loop through optionally configured paths
	for _, path := range configPaths {
		viper.AddConfigPath(path)
	}

	// set up list of alternate config file paths
	viper.AddConfigPath(c.homeConfigDirectory())
	viper.AddConfigPath("/usr/local/etc")
	viper.AddConfigPath(".")

	populateConfig(c)
}
