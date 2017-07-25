package main

import (
	"log"
	"os"
	"os/signal"

	"dropboy"
	"dropboy/config"
)

// var done = make(chan bool, 2)

func main() {
	log.Println("Starting Dropboy.")

	dboy := dropboy.NewDropboy()

	c := new(config.Config)
	c.Configure()
	dboy.RegisterWatchesFromConfig(c)

	irq := make(chan os.Signal, 1)
	signal.Notify(irq, os.Interrupt)
	go func() {
		for range irq {
			log.Println("Stopping Dropboy.")
			dboy.Stop()
			os.Exit(1)
		}
	}()

	for {
		dboy.HandleFilesystemEvents(dboy.Notifier.Events)
	}

}
