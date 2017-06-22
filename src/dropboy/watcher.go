package dropboy

import (
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"

	"dropboy/config"
)

// TODO: rename this whole thing to Dropboy -- it's really the top level thing

type watcher struct {
	Watches  map[string][]string
	Notifier *fsnotify.Watcher
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

	// TODO: I don't know what this is really doing for us here
	// a list of paths for later use?  The filesystem watcher is stateless.
	w.Watches[path] = actions

	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}

	err = w.Notifier.Add(absPath)
	if err != nil {
		log.Fatal(err)
	}
}

func (w *watcher) RegisterWatchesFromConfig(config config.Config) {
	for _, watch := range config.Watches {
		actions := []string{}
		for _, action := range watch.Actions {
			actions = append(actions, action.Type)
		}
		w.Register(watch.Path, actions)
	}
}
