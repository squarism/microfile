package config

import (
	"fmt"
	"log"
	"os"
	"testing"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
)

// TODO: if relative directory paths work here then we should just replace
// this test helper with a string like
// var directoryToWatch string = "./../../fixtures"
func fixturesDirectory() string {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%s/../../../fixtures", pwd)
}

func patchEnv(key, value string) func() {
	bck := os.Getenv(key)
	deferFunc := func() {
		os.Setenv(key, bck)
	}

	os.Setenv(key, value)
	return deferFunc
}

func TestDefaultURLConfig(t *testing.T) {
	c := new(Config)
	c.Configure(fixturesDirectory())

	assert.Equal(t, "http://localhost:9876", c.DefaultURL, "Test config defaults")
}

// This deeply tests the config in fixtures
func TestActions(t *testing.T) {
	var expected = []Action{
		{
			Type: "http",
			Options: map[string]string{
				"send_file": "true",
			},
		},
	}

	c := new(Config)
	c.Configure(fixturesDirectory())
	watches := c.Watches

	assert.Equal(t, expected, watches[0].Actions, "Test actions from config file")
}

func TestHomeConfigDir(t *testing.T) {
	// overtesting homedir probably
	// we have to setup the test here,
	// disable home directory caching and insert ENV stuff
	homedir.DisableCache = true
	defer patchEnv("HOME", "/custom/path")()

	c := new(Config)
	assert.Equal(t, "/custom/path/.microfile/", c.homeConfigDirectory())
}
