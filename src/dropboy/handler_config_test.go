package dropboy

import (
	"testing"

	"github.com/fsnotify/fsnotify"
	"github.com/stretchr/testify/assert"

	"dropboy/config"
	"dropboy/handler"
)

var validConfig = config.Config{
	DefaultURL: "http://localhost:3000/",
	Watches: []config.Watch{
		{
			Path: "/tmp/foo",
			Actions: []config.Action{
				{
					Type:    "http",
					Options: map[string]string{"path": "/remote/server/path"},
				},
			},
		},
	},
}

var invalidActionConfig = config.Config{
	Watches: []config.Watch{
		{
			Path: "/tmp/foo",
			Actions: []config.Action{
				{
					Type: "telnet",
				},
			},
		},
	},
}

func TestValidateValidHandlers(t *testing.T) {
	handlerConfig := new(HandlerConfig)
	err := handlerConfig.Validate(validConfig)

	assert.Nil(t, err)
}

func TestValidateInValidHandlers(t *testing.T) {
	handlerConfig := new(HandlerConfig)
	err := handlerConfig.Validate(invalidActionConfig)

	assert.NotNil(t, err)
}

func TestHandlersFor(t *testing.T) {
	handlerConfig := new(HandlerConfig)
	handlers := handlerConfig.HandlersFor("/tmp/foo/bleh.txt", validConfig)

	path := handlers[0].(*handler.HTTP).Path

	assert.Equal(t, "/remote/server/path", path)
}

func TestDefaultsFromHandlerConfig(t *testing.T) {
	handlerConfig := new(HandlerConfig)
	handlers := handlerConfig.HandlersFor("/tmp/foo/bleh.txt", validConfig)

	expected := "http://localhost:3000/"
	url := handlers[0].(*handler.HTTP).DefaultURL

	assert.Equal(t, expected, url)
}

func TestIgnoreEvents(t *testing.T) {
	event := fsnotify.Event{Name: "bleh.txt", Op: fsnotify.Chmod}
	handlerConfig := new(HandlerConfig)

	result := handlerConfig.IsRelevantEvent(event)

	assert.Equal(t, false, result)
}
