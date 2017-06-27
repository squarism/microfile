package dropboy

import (
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"

	"dropboy/config"
)

type dropboy struct {
	Watches       map[string][]string
	Notifier      *fsnotify.Watcher
	HandlerConfig HandlerFinder
	Config        *config.Config
}

func NewDropboy() dropboy {
	notifier, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	d := dropboy{}
	d.Notifier = notifier

	return d
}

func (w *dropboy) Stop() {
	w.Notifier.Close()
}

func (w *dropboy) Register(path string, actions []string) {
	if w.Watches == nil {
		w.Watches = make(map[string][]string)
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}

	err = w.Notifier.Add(absPath)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: I don't know what this is really doing for us here
	// a list of paths for later use?  The filesystem dropboy is stateless.
	w.Watches[absPath] = actions
}

func (w *dropboy) RegisterWatchesFromConfig(c *config.Config) {
	w.Config = c
	w.HandlerConfig = &HandlerConfig{}

	for _, watch := range c.Watches {
		actions := []string{}
		for _, action := range watch.Actions {
			actions = append(actions, action.Type)
		}
		w.Register(watch.Path, actions)
	}
}

// TODO: we are completely ignoring the dropboy.Errors channel here
// create another method to handle those
func (w *dropboy) HandleFilesystemEvents(channel chan fsnotify.Event) {
	select {
	case event := <-channel:
		path := event.Name
		handlers := w.HandlerConfig.HandlersFor(path, *w.Config)

		// We could reflect and determine the handler type here and switch on it
		// for custom behavior outside the Handle method but why?  Why not just
		// send the event to the method defined by the interface?  We validate the config elsewhere.
		// We shouldn't have illegal handlers.
		for _, h := range handlers {
			h.Handle(event)
		}
	}
}
