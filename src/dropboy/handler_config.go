package dropboy

import (
	"fmt"
	"path/filepath"

	"dropboy/config"
	"dropboy/handler"
)

var handlerNames = [2]string{"http", "log"}

var handlerBleh = map[string]handler.Handler{
	"http": &handler.HTTP{},
	"log":  &handler.Log{},
}

type HandlerConfig struct {
}

func (handlerConfig *HandlerConfig) Validate(config config.Config) error {
	numberOfInvalidActions := 0

	for _, w := range config.Watches {
		for _, a := range w.Actions {
			if validHandlerName(a.Type) == false {
				numberOfInvalidActions++
			}
		}
	}

	if numberOfInvalidActions > 0 {
		return fmt.Errorf("Invalid action detected in config.  Valid actions are %v", handlerNames)
	}

	return nil
}

func (handlerConfig *HandlerConfig) HandlersFor(path string, config config.Config) []handler.Handler {
	// why do we need to do this?  Well, we're probably going to be receiving
	// a path from a filesystem event like /tmp/foo/file.txt but we are watching
	// /tmp/foo right?  So we need to be looking up the handler for the directory basename.
	// TODO: this is probably going to be a huge headache when we go recursive?  Or maybe not
	// if we recurse ourselves and just register all the child paths as watches?
	dir := filepath.Dir(path)

	handlers := []handler.Handler{}
	for _, watch := range config.Watches {
		if watch.Path == dir {
			for _, action := range watch.Actions {
				switch action.Type {
				case "http":
					handler := &handler.HTTP{}
					handler.Init(action)
					handlers = append(handlers, handler)
				}
			}
		}
	}

	return handlers
}

func validHandlerName(name string) bool {
	for _, n := range handlerNames {
		if n == name {
			return true
		}
	}
	return false
}