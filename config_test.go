package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultURL(t *testing.T) {
	config := Config{}
	c := config.Load()

	assert.Equal(t, "http://localhost:3000", c.DefaultUrl, "Test config")
}

func TestTriggers(t *testing.T) {
	var expected = map[string][]string{
		"/var/www/resumes/dropbox": {"/api/convert_doc"},
	}

	config := Config{}
	c := config.Load()

	assert.Equal(t, expected, c.Triggers, "Test config")
}

func TestAlternateFile(t *testing.T) {
	config := Config{}
	c := config.Load("./dropboy_test.yml")

	assert.Equal(t, "alternate file", c.DefaultUrl, "Test config")
}
