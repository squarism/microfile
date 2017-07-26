package config

import (
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/hcl"
	"github.com/mitchellh/go-homedir" // avoids cgo cross compile issues
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper" // leaning heavily on viper for configuration
)

type Config struct {
	DefaultURL string  `hcl:"default_url"`
	LogFile    string  `hcl:"log_file"`
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
		log.WithFields(log.Fields{"lifecycle": "config"}).Fatal("Can't read config file")
	}

	// because viper doesn't pass HCL flags down for auto-key niceness
	// lets just bypass viper a bit
	viperConfigFileContents, err := ioutil.ReadFile(viper.GetViper().ConfigFileUsed())
	if err != nil {
		log.WithFields(log.Fields{"lifecycle": "config"}).Fatal("Can't read config file contents")
	}

	// read everything into the config reference (c), this is like &config from the HCL examples
	err = hcl.Decode(c, string(viperConfigFileContents))
	if err != nil {
		log.WithFields(log.Fields{"lifecycle": "config"}).Fatal("Problem with the HCL config file")
	}
}

func (c *Config) homeConfigDirectory() string {
	home, err := homedir.Dir()
	if err != nil {
		log.WithFields(log.Fields{"lifecycle": "config"}).Fatal("Cannot tell what the home directory is")
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

	// TODO - add ENV option for the config path

	// set up list of alternate config file paths
	viper.AddConfigPath(c.homeConfigDirectory())
	viper.AddConfigPath("/usr/local/etc")
	viper.AddConfigPath(".")

	populateConfig(c)
}
