package dropboy

import (
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"

	"dropboy/config"
)

// TODO: rename this whole thing to Dropboy -- it's really the top level thing

type watcher struct {
	Watches       map[string][]string
	Notifier      *fsnotify.Watcher
	HandlerConfig HandlerFinder
	Config        *config.Config
}

func NewWatcher() watcher {
	notifier, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	w := watcher{}
	w.Notifier = notifier

	return w
}

func (w *watcher) Stop() {
	w.Notifier.Close()
}

func (w *watcher) Register(path string, actions []string) {
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
	// a list of paths for later use?  The filesystem watcher is stateless.
	w.Watches[absPath] = actions
}

func (w *watcher) RegisterWatchesFromConfig(c *config.Config) {
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

// TODO: we are completely ignoring the watcher.Errors channel here
// create another method to handle those
func (w *watcher) HandleFilesystemEvents(channel chan fsnotify.Event) {
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
