package dropboy

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

var workDir string = "/var/dropboy"

func TestHashFile(t *testing.T) {
	locker := NewLocker(workDir)
	filename := "/var/www/uploads/song.mp3"

	result := locker.Hash(filename)

	expected := "733388170efcc73067c22d5e98c55008"
	assert.Equal(t, expected, result)
}

func TestLockFilePushesToList(t *testing.T) {
	locker := NewLocker(workDir)
	filename := "/var/www/uploads/song.mp3"
	hash := "733388170efcc73067c22d5e98c55008"

	locker.Lock(filename, hash)

	assert.Len(t, locker.AllActive, 1)
}

func TestEnsureWorkDirectory(t *testing.T) {
	locker := NewLocker(workDir)
	locker.AppFs = afero.NewMemMapFs() // for testing

	locker.ensureWorkDirectory(workDir)
	result, _ := afero.Exists(locker.AppFs, workDir)

	expected := true
	assert.Equal(t, expected, result)
}

func TestHashFileName(t *testing.T) {
	locker := NewLocker(workDir)
	hash := "733388170efcc73067c22d5e98c55008"

	result := locker.hashFilename(hash)

	expected := "/var/dropboy/733388170efcc73067c22d5e98c55008"
	assert.Equal(t, expected, result)
}
