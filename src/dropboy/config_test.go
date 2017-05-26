package dropboy

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func fixturesDirectory() string {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%s/../../fixtures", pwd)
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

	assert.Equal(t, "http://localhost:9876", c.DefaultUrl, "Test config")
}

func TestTriggers(t *testing.T) {
	var expected = map[string][]string{
		"/var/www/resumes/dropbox": {"/api/convert_doc"},
	}

	c := new(Config)
	c.Configure(fixturesDirectory())

	assert.Equal(t, expected, c.Triggers, "Test triggers")
}

func TestHomeConfigDir(t *testing.T) {
	// overtesting homedir probably
	// we have to setup the test here,
	// disable home directory caching and insert ENV stuff
	homedir.DisableCache = true
	defer patchEnv("HOME", "/custom/path")()

	c := new(Config)
	assert.Equal(t, "/custom/path/.dropboy/", c.homeConfigDirectory())
}
