package microfile

import (
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"

	"microfile/config"
)

type microfile struct {
	Watches       map[string][]string
	Notifier      *fsnotify.Watcher
	HandlerConfig HandlerFinder
	Config        *config.Config
	Locker        locker
}

func NewMicrofile(config *config.Config) microfile {
	notifier, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	mf := microfile{}
	mf.Watches = make(map[string][]string) // init a nil map
	mf.Notifier = notifier
	mf.Locker = NewLocker()

	mf.LoadConfig(config)
	return mf
}

func (m *microfile) Stop() {
	m.Notifier.Close()
}

func (m *microfile) Register(path string, actions []string) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}

	err = m.Notifier.Add(absPath)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: I don't know what this is really doing for us here
	// a list of paths for later use?  The filesystem microfile is stateless.
	m.Watches[absPath] = actions
}

func (m *microfile) LoadConfig(c *config.Config) {
	m.Config = c
	m.HandlerConfig = &HandlerConfig{}

	m.setupWorkDirectory()

	for _, watch := range c.Watches {
		actions := []string{}
		for _, action := range watch.Actions {
			actions = append(actions, action.Type)
		}
		m.Register(watch.Path, actions)
	}
}

// TODO: we are completely ignoring the microfile.Errors channel here
// create another method to handle those
func (m *microfile) HandleFilesystemEvents(channel chan fsnotify.Event) {
	select {
	case event := <-channel:
		path := event.Name
		handlers := m.HandlerConfig.HandlersFor(path, *m.Config)

		// We could reflect and determine the handler type here and switch on it
		// for custom behavior outside the Handle method but why?  Why not just
		// send the event to the method defined by the interface?  We validate the config elsewhere.
		// We shouldn't have illegal handlers.
		for _, h := range handlers {
			if m.isRelevantEvent(event) {
				h.Handle(event)
			}
		}
	}
}

// Handles overriding the working directory for file locks from the config
func (m *microfile) setupWorkDirectory() {
	if m.Config.WorkDirectory != "" {
		m.Locker.WorkDirectory = m.Config.WorkDirectory
	} else {
		log.Debug("Using default work directory")
	}
}

// global ignore of sorts
func (m *microfile) isRelevantEvent(event fsnotify.Event) bool {
	return (event.Op != fsnotify.Chmod)
}
