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
	watcher.Register("/tmp", []string{"http://localhost:3000/tmp"})

	assert.Equal(t, len(watcher.Watches), 1)
}

func TestRegisterWatches(t *testing.T) {
	watcher := Watcher{}
	watcher.Register("/tmp", []string{"http://localhost:3000/tmp"})
	watcher.Register("/var/tmp", []string{"http://localhost:3000/var"})

	assert.Equal(t, len(watcher.Watches), 2)
}
