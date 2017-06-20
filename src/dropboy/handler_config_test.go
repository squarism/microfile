package dropboy

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"dropboy/config"
	"dropboy/handler"
)

var validConfig = config.Config{
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
