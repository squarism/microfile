package dropboy

import (
	"fmt"
	"github.com/mitchellh/go-homedir" // avoids cgo cross compile issues
	"github.com/spf13/viper"          // leaning heavily on viper for configuration
	"log"
)

// Config represents a ringu configuration
type Config struct {
	DefaultUrl string // `yaml:"default_url"`
	Triggers   map[string][]string
}

// read in config defaults and populate our config struct
func populateConfig(c *Config) {
	err := viper.ReadInConfig()

	c.DefaultUrl = viper.GetString("default_url")
	c.Triggers = viper.GetStringMapStringSlice("triggers")

	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func (c *Config) homeConfigDirectory() string {
	home, err := homedir.Dir()
	if err != nil {
		log.Panic(err)
	}

	// This will allow users to put the config in ~/.ringu/ringu.yml
	homeConfigDirectory := fmt.Sprintf("%s/.dropboy/", home)
	return homeConfigDirectory
}

// set config options and default values in the config
func setConfigurationDefaults() {
	viper.SetDefault("DefaultUrl", "http://localhost:3000")
	viper.SetConfigName("dropboy")
	viper.SetConfigType("yaml")
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
