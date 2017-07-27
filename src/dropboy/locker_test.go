package dropboy

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

var filename string = "/var/www/uploads/song.mp3"

func TestHashFile(t *testing.T) {
	locker := NewLocker()

	result := locker.Hash(filename)

	expected := "733388170efcc73067c22d5e98c55008"
	assert.Equal(t, expected, result)
}

func TestLockFilePushesToList(t *testing.T) {
	locker := NewLocker()
	locker.AppFs = afero.NewMemMapFs()
	locker.AppFs.MkdirAll(filepath.Base(filename), 0755)
	locker.AppFs.Create(filename)

	locker.Lock(filename)

	assert.Len(t, locker.AllActive, 1)
}

func TestEnsureWorkDirectory(t *testing.T) {
	locker := NewLocker()
	workDir := "/1/2/3"
	locker.AppFs = afero.NewMemMapFs() // for testing

	locker.ensureWorkDirectory(workDir)
	result, _ := afero.DirExists(locker.AppFs, workDir)

	expected := true
	assert.Equal(t, expected, result)
}

func TestHashFileName(t *testing.T) {
	locker := NewLocker()
	locker.WorkDirectory = "/var/dropboy"
	hash := "733388170efcc73067c22d5e98c55008"

	result := locker.hashFilename(hash)

	expected := "/var/dropboy/733388170efcc73067c22d5e98c55008"
	assert.Equal(t, expected, result)
}

func TestLocking(t *testing.T) {
	locker := NewLocker()
	locker.WorkDirectory = "/var/dropboy"

	locker.AppFs = afero.NewMemMapFs()
	locker.AppFs.MkdirAll(filepath.Base(filename), 0755)
	locker.AppFs.Create(filename)

	locker.Lock(filename)
	expectedLockedFile := "/var/dropboy/733388170efcc73067c22d5e98c55008"

	result, _ := afero.Exists(locker.AppFs, expectedLockedFile)
	expected := true
	assert.Equal(t, expected, result)
}