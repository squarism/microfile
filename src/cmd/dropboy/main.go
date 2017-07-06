package main

import (
	"log"

	"dropboy"
	"dropboy/config"
)

// var done = make(chan bool, 2)

func main() {
	log.Println("Starting Dropboy ...")

	dboy := dropboy.NewDropboy()

	c := new(config.Config)
	c.Configure()
	dboy.RegisterWatchesFromConfig(c)

	dboy.HandleFilesystemEvents(dboy.Notifier.Events)
	dboy.Stop()

	log.Println("Stopping Dropboy ...")
}
