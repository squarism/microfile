package dropboy

import (
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	Watches map[string][]string
}

func (w *Watcher) Register(watcher *fsnotify.Watcher, path string, urls []string) {
	if w.Watches == nil {
		w.Watches = make(map[string][]string)
	}

	w.Watches[path] = urls

	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}

	err = watcher.Add(absPath)
	if err != nil {
		log.Fatal(err)
	}
}
