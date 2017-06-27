package handler

import (
	"log"

	"github.com/fsnotify/fsnotify"

	"dropboy/config"
)

type HTTP struct {
	Path string
}

func (http HTTP) Handle(event fsnotify.Event) {
	log.Println("reacting with a HTTP handler: ", event)
}

func (http *HTTP) Init(action config.Action) error {
	http.Path = action.Options["path"]
	return nil
}
