package handler

import (
	"log"

	"github.com/fsnotify/fsnotify"

	"dropboy/config"
)

type Log struct {
}

func (logger Log) Handle(event fsnotify.Event) {
	log.Println("event:", event)
}

func (logger *Log) Init(action config.Action) error {
	return nil
}
