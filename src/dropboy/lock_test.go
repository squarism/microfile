package dropboy

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestHashFile(t *testing.T) {
	locker := NewLocker()
	filename := "/var/www/uploads/song.mp3"

	expected := "733388170efcc73067c22d5e98c55008"
	hash := locker.Hash(filename)

	assert.Equal(t, expected, hash)
}

func TestLockFilePushesToList(t *testing.T) {
	locker := NewLocker()
	filename := "/var/www/uploads/song.mp3"
	hash := "733388170efcc73067c22d5e98c55008"

	locker.Lock(filename, hash)

	assert.Len(t, locker.AllActive, 1)
}

func TestLockSkipsPushWhenAlreadyThere(t *testing.T) {
	locker := NewLocker()
	filename := "/var/www/uploads/song.mp3"
	hash := "733388170efcc73067c22d5e98c55008"
	locker.AllActive[filename] = hash

	err := locker.Lock(filename, hash)

	assert.EqualError(t, err, "Already processing /var/www/uploads/song.mp3")
}

func TestEnsureWorkDirectory(t *testing.T) {
	workDir := "/var/dropboy"
	locker := NewLocker()
	locker.AppFs = afero.NewMemMapFs() // for testing

	locker.EnsureWorkDirectory(workDir)
	exists, _ := afero.Exists(locker.AppFs, workDir)

	assert.True(t, exists)
}
