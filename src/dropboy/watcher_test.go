package dropboy

import (
	"fmt"
	"log"
	"path/filepath"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"dropboy/config"
	"dropboy/handler"
)

func fixturesDirectory() string {
	d, _ := filepath.Abs("./../../fixtures/dropbox")
	return d
}

var directoryToWatch string = fixturesDirectory()

var watcherValidConfig = config.Config{
	DefaultURL: "http://somewhere/valid_config",
	Watches: []config.Watch{
		{
			Path: directoryToWatch,
			Actions: []config.Action{
				{
					Type:    "http",
					Options: map[string]string{"path": "/remote/server/path"},
				},
			},
		},
	},
}

func TestRegisterNothing(t *testing.T) {
	watcher := NewWatcher()
	defer watcher.Stop()

	assert.Equal(t, 0, len(watcher.Watches))
}

func TestRegisterWatch(t *testing.T) {
	watcher := NewWatcher()
	defer watcher.Stop()

	watcher.Register(directoryToWatch, []string{"http://localhost:3000/resumes"})

	assert.Equal(t, 1, len(watcher.Watches))
}

func TestRegisterWatches(t *testing.T) {
	realEstateDirectory := fmt.Sprint(directoryToWatch, "/real_estate")
	musicDirectory := fmt.Sprint(directoryToWatch, "/music")
	watcher := NewWatcher()
	defer watcher.Stop()

	watcher.Register(realEstateDirectory, []string{"http://localhost:3000/image_shrinker"})
	watcher.Register(musicDirectory, []string{
		"http://localhost:3000/copyright_alerter",
		"http://localhost:3000/recompressor",
	})

	assert.Equal(t, 2, len(watcher.Watches))
}

func TestRegisterFromConfig(t *testing.T) {
	dir, err := filepath.Abs(directoryToWatch)
	if err != nil {
		log.Fatal("Fixtures directory is missing.")
	}

	watcher := NewWatcher()
	defer watcher.Stop()
	expected := map[string][]string{
		dir: []string{"http"},
	}

	watcher.RegisterWatchesFromConfig(&watcherValidConfig)

	assert.Equal(t, expected, watcher.Watches)
}

func TestRememberConfig(t *testing.T) {
	watcher := NewWatcher()
	watcher.RegisterWatchesFromConfig(&watcherValidConfig)
	expected := "http://somewhere/valid_config"

	assert.Equal(t, expected, watcher.Config.DefaultURL)
}

// wow mocking in Go is crazy hard
type MockHandlerConfig struct {
	mock.Mock
}

func (m *MockHandlerConfig) HandlersFor(path string, config config.Config) []handler.Handler {
	m.Called(path)
	return mockHandlers
}

type MockHandler struct {
	mock.Mock
}

func (m *MockHandler) Handle(event fsnotify.Event) {
	log.Println("OMG MOCK HANDLER")
	m.Called(event)
}

func (m *MockHandler) Init(action config.Action) error {
	return nil
}

var mockHandlers = []handler.Handler{}

func TestIncomingFilesystemEvents(t *testing.T) {
	watcher := NewWatcher()
	watcher.RegisterWatchesFromConfig(&watcherValidConfig)
	mockHandlerConfig := new(MockHandlerConfig)
	changedFile, _ := filepath.Abs(fmt.Sprintf("%s/file_that_changed.txt", directoryToWatch))

	channel := make(chan fsnotify.Event)
	event := fsnotify.Event{Name: changedFile}

	mockHandler := new(MockHandler)
	mockHandler.On("Handle", event)
	mockHandlers = append(mockHandlers, mockHandler)

	watcher.HandlerConfig = mockHandlerConfig
	mockHandlerConfig.On("HandlersFor", changedFile).Return(mockHandlers)

	done := make(chan bool, 2)

	go func(done chan bool) {
		watcher.HandleFilesystemEvents(channel)
		done <- true
	}(done)

	go func(done chan bool) {
		channel <- event
		time.Sleep(time.Millisecond * 1000)
		done <- true
	}(done)

	<-done

	mockHandlerConfig.AssertExpectations(t)
}
