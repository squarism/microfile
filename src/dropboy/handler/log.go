package handler

import (
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"

	"dropboy/config"
)

type Log struct {
}

func (logger Log) Handle(event fsnotify.Event) {
	log.WithFields(
		log.Fields{
			"handler":  "log",
			"filename": event.Name,
			"event":    event.Op.String(),
		}).Info("Filesystem Event")
}

func (logger *Log) Init(action config.Action) error {
	return nil
}
