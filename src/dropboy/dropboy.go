package dropboy

import (
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"

	"dropboy/config"
)

type dropboy struct {
	Watches       map[string][]string
	Notifier      *fsnotify.Watcher
	HandlerConfig HandlerFinder
	Config        *config.Config
	Locker        locker
}

func NewDropboy() dropboy {
	notifier, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	d := dropboy{}
	d.Notifier = notifier
	d.Locker = NewLocker()

	return d
}

func (d *dropboy) Stop() {
	d.Notifier.Close()
}

func (d *dropboy) Register(path string, actions []string) {
	if d.Watches == nil {
		d.Watches = make(map[string][]string)
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}

	err = d.Notifier.Add(absPath)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: I don't know what this is really doing for us here
	// a list of paths for later use?  The filesystem dropboy is stateless.
	d.Watches[absPath] = actions
}

func (d *dropboy) LoadConfig(c *config.Config) {
	d.Config = c
	d.HandlerConfig = &HandlerConfig{}

	d.setupWorkDirectory()

	for _, watch := range c.Watches {
		actions := []string{}
		for _, action := range watch.Actions {
			actions = append(actions, action.Type)
		}
		d.Register(watch.Path, actions)
	}
}

// TODO: we are completely ignoring the dropboy.Errors channel here
// create another method to handle those
func (d *dropboy) HandleFilesystemEvents(channel chan fsnotify.Event) {
	select {
	case event := <-channel:
		path := event.Name
		handlers := d.HandlerConfig.HandlersFor(path, *d.Config)

		// We could reflect and determine the handler type here and switch on it
		// for custom behavior outside the Handle method but why?  Why not just
		// send the event to the method defined by the interface?  We validate the config elsewhere.
		// We shouldn't have illegal handlers.
		for _, h := range handlers {
			if d.isRelevantEvent(event) {
				h.Handle(event)
			}
		}
	}
}

// Handles overriding the working directory for file locks from the config
func (d *dropboy) setupWorkDirectory() {
	// spew.Dump(d.Config)

	if d.Config.WorkDirectory != "" {
		d.Locker.WorkDirectory = d.Config.WorkDirectory
	} else {
		log.Info("Using default work directory")
	}
}

// global ignore of sorts
func (d *dropboy) isRelevantEvent(event fsnotify.Event) bool {
	return (event.Op != fsnotify.Chmod)
}
