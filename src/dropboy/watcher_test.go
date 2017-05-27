package dropboy

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegisterNothing(t *testing.T) {
	watcher := Watcher{}

	assert.Equal(t, len(watcher.Watches), 0)
}

func TestRegisterWatch(t *testing.T) {
	watcher := Watcher{}
	watcher.Register("/resumes", []string{"http://localhost:3000/resumes"})

	assert.Equal(t, len(watcher.Watches), 1)
}

func TestRegisterWatches(t *testing.T) {
	watcher := Watcher{}
	watcher.Register("/real_estate", []string{"http://localhost:3000/image_shrinker"})
	watcher.Register("/music", []string{
		"http://localhost:3000/copyright_alerter",
		"http://localhost:3000/recompressor",
	})

	assert.Equal(t, len(watcher.Watches), 2)
}
