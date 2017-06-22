package dropboy

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"dropboy/config"
)

var directoryToWatch string = "./../../fixtures/dropbox"

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
	var config = config.Config{
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

	watcher := NewWatcher()
	defer watcher.Stop()
	expected := map[string][]string{
		directoryToWatch: []string{"http"},
	}

	watcher.RegisterWatchesFromConfig(config)
	assert.Equal(t, expected, watcher.Watches)
}
