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

var dropboyValidConfig = config.Config{
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
	dropboy := NewDropboy()
	defer dropboy.Stop()

	assert.Equal(t, 0, len(dropboy.Watches))
}

func TestRegisterWatch(t *testing.T) {
	dropboy := NewDropboy()
	defer dropboy.Stop()

	dropboy.Register(directoryToWatch, []string{"http://localhost:3000/resumes"})

	assert.Equal(t, 1, len(dropboy.Watches))
}

func TestRegisterWatches(t *testing.T) {
	realEstateDirectory := fmt.Sprint(directoryToWatch, "/real_estate")
	musicDirectory := fmt.Sprint(directoryToWatch, "/music")
	dropboy := NewDropboy()
	defer dropboy.Stop()

	dropboy.Register(realEstateDirectory, []string{"http://localhost:3000/image_shrinker"})
	dropboy.Register(musicDirectory, []string{
		"http://localhost:3000/copyright_alerter",
		"http://localhost:3000/recompressor",
	})

	assert.Equal(t, 2, len(dropboy.Watches))
}

func TestRegisterFromConfig(t *testing.T) {
	dir, err := filepath.Abs(directoryToWatch)
	if err != nil {
		log.Fatal("Fixtures directory is missing.")
	}

	dropboy := NewDropboy()
	defer dropboy.Stop()
	expected := map[string][]string{
		dir: []string{"http"},
	}

	dropboy.RegisterWatchesFromConfig(&dropboyValidConfig)

	assert.Equal(t, expected, dropboy.Watches)
}

func TestRememberConfig(t *testing.T) {
	dropboy := NewDropboy()
	dropboy.RegisterWatchesFromConfig(&dropboyValidConfig)
	expected := "http://somewhere/valid_config"

	assert.Equal(t, expected, dropboy.Config.DefaultURL)
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
	m.Called(event)
}

func (m *MockHandler) Init(action config.Action) error {
	return nil
}

var mockHandlers = []handler.Handler{}

func TestIncomingFilesystemEvents(t *testing.T) {
	dropboy := NewDropboy()
	dropboy.RegisterWatchesFromConfig(&dropboyValidConfig)
	mockHandlerConfig := new(MockHandlerConfig)
	changedFile, _ := filepath.Abs(fmt.Sprintf("%s/file_that_changed.txt", directoryToWatch))

	channel := make(chan fsnotify.Event)
	event := fsnotify.Event{Name: changedFile}

	mockHandler := new(MockHandler)
	mockHandler.On("Handle", event)
	mockHandlers = append(mockHandlers, mockHandler)

	dropboy.HandlerConfig = mockHandlerConfig
	mockHandlerConfig.On("HandlersFor", changedFile).Return(mockHandlers)

	// we need to make sure our test doesn't exit too soon
	done := make(chan bool, 2)

	go func(done chan bool) {
		dropboy.HandleFilesystemEvents(channel)
		done <- true
	}(done)

	go func(done chan bool) {
		channel <- event
		time.Sleep(time.Millisecond)
		done <- true
	}(done)

	<-done

	mockHandlerConfig.AssertExpectations(t)
}

func TestIgnoreEvents(t *testing.T) {
	event := fsnotify.Event{Name: "bleh.txt", Op: fsnotify.Chmod}
	dropboy := NewDropboy()

	result := dropboy.isRelevantEvent(event)

	assert.Equal(t, false, result)
}
