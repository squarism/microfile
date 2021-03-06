package main

import (
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"

	"microfile"
	"microfile/config"
)

// This is the main daemon binary for Microfile.
// It handles entry into other parts of the programs and tries not to
// do work itself that is unrelated to the CLI, invoking or the daemon itself.
func main() {
	c := new(config.Config)
	c.Configure()

	mf := microfile.NewMicrofile(c)

	setupLogging(c)

	log.WithFields(log.Fields{"lifecycle": "startup"}).Info("Starting Microfile")

	// Handling Ctrl-C as an interrupt (irq)
	irq := make(chan os.Signal, 1)
	signal.Notify(irq, os.Interrupt)
	go func() {
		for range irq {
			log.WithFields(log.Fields{"lifecycle": "shutdown"}).Info("Stopping Microfile")
			mf.Stop()
			os.Exit(1)
		}
	}()

	// Infinite loop until something interrupts
	for {
		mf.HandleFilesystemEvents(mf.Notifier.Events)
	}
}

func setupLogging(c *config.Config) {
	// you can `DEBUG=true ./bin/microfile` to show more messages
	if os.Getenv("DEBUG") != "" {
		// show debug messages through Logrus
		log.SetLevel(log.DebugLevel)
	}

	if c.LogFile != "" {
		logFile, err := os.OpenFile(c.LogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.WithFields(log.Fields{"file": "c.LogFile", "lifecycle": "startup"}).Fatal("Can't open log file")
		}
		log.SetOutput(logFile)
	} else {
		log.Debug("Using default log")
	}
}
