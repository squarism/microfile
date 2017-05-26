package main

type Watcher struct {
	Watches map[string][]string
}

func (w *Watcher) Register(path string, urls []string) {
	if w.Watches == nil {
		w.Watches = make(map[string][]string)
	}

	w.Watches[path] = urls
}
