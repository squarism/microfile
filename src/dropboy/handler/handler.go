package handler

import (
	"github.com/fsnotify/fsnotify"

	"dropboy/config"
)

type Handler interface {
	Handle(event fsnotify.Event)
	Init(action config.Action) error
}
