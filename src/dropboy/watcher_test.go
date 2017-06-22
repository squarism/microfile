package dropboy

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"dropboy/config"
)

var directoryToWatch string = "./../../fixtures/dropbox"

// We have to duplicate this testing helper method from another test file.
// Shared test code in Go is either duplicated or documented from godoc.  Which do I want?
func fixturesDirectory() string {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%s/../../fixtures", pwd)
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
	publicFixturesPath := fmt.Sprintf("%s/dropbox", fixturesDirectory())

	var config = config.Config{
		Watches: []config.Watch{
			{
				Path: publicFixturesPath,
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
		publicFixturesPath: []string{"http"},
	}

	watcher.RegisterWatchesFromConfig(config)
	assert.Equal(t, expected, watcher.Watches)
}
