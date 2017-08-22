package microfile

import (
	"fmt"
	"log"
	"path/filepath"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"microfile/config"
	"microfile/handler"
)

func fixturesDirectory() string {
	directory, _ := filepath.Abs("./../../fixtures/dropbox")
	return directory
}

var directoryToWatch string = fixturesDirectory()

var microfileValidConfig = config.Config{
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

var microfileEmptyConfig = config.Config{
	DefaultURL: "http://somewhere/valid_config",
	Watches:    []config.Watch{},
}

func TestRegisterNothing(t *testing.T) {
	microfile := NewMicrofile(&microfileEmptyConfig)

	assert.Equal(t, 0, len(microfile.Watches))
}

func TestRegisterWatch(t *testing.T) {
	microfile := NewMicrofile(&microfileEmptyConfig)

	microfile.Register(directoryToWatch, []string{"http://localhost:3000/resumes"})

	assert.Equal(t, 1, len(microfile.Watches))
}

func TestRegisterWatches(t *testing.T) {
	realEstateDirectory := fmt.Sprint(directoryToWatch, "/real_estate")
	musicDirectory := fmt.Sprint(directoryToWatch, "/music")
	microfile := NewMicrofile(&microfileEmptyConfig)

	microfile.Register(realEstateDirectory, []string{"http://localhost:3000/image_shrinker"})
	microfile.Register(musicDirectory, []string{
		"http://localhost:3000/copyright_alerter",
		"http://localhost:3000/recompressor",
	})

	assert.Equal(t, 2, len(microfile.Watches))
}

func TestRegisterFromConfig(t *testing.T) {
	dir, err := filepath.Abs(directoryToWatch)
	if err != nil {
		log.Fatal("Fixtures directory is missing.")
	}

	microfile := NewMicrofile(&microfileValidConfig)

	expected := map[string][]string{
		dir: []string{"http"},
	}

	assert.Equal(t, expected, microfile.Watches)
}

func TestRememberConfig(t *testing.T) {
	microfile := NewMicrofile(&microfileValidConfig)
	expected := "http://somewhere/valid_config"

	assert.Equal(t, expected, microfile.Config.DefaultURL)
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
	microfile := NewMicrofile(&microfileValidConfig)
	mockHandlerConfig := new(MockHandlerConfig)
	changedFile, _ := filepath.Abs(fmt.Sprintf("%s/file_that_changed.txt", directoryToWatch))

	channel := make(chan fsnotify.Event)
	event := fsnotify.Event{Name: changedFile}

	mockHandler := new(MockHandler)
	mockHandler.On("Handle", event)
	mockHandlers = append(mockHandlers, mockHandler)

	microfile.HandlerConfig = mockHandlerConfig
	mockHandlerConfig.On("HandlersFor", changedFile).Return(mockHandlers)

	// we need to make sure our test doesn't exit too soon
	done := make(chan bool, 2)

	go func(done chan bool) {
		microfile.HandleFilesystemEvents(channel)
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
	microfile := NewMicrofile(&microfileValidConfig)

	result := microfile.isRelevantEvent(event)

	assert.Equal(t, false, result)
}
