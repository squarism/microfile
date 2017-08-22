package dropboy

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type locker struct {
	AllActive     map[string]string
	AppFs         afero.Fs
	WorkDirectory string // TODO: this should be able to be overridden in the config
}

func NewLocker() locker {
	l := locker{}
	l.AllActive = make(map[string]string)
	l.AppFs = afero.NewOsFs()
	l.WorkDirectory = filepath.Join(os.TempDir(), "dropboy")

	return l
}

func (l *locker) Hash(filename string) string {
	hasher := md5.New()
	io.WriteString(hasher, filename)
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

func (l *locker) Lock(filename string) {
	if _, ok := l.AllActive[filename]; ok {
		log.WithFields(log.Fields{
			"workDir":  l.WorkDirectory,
			"filename": filename},
		).Info("Already processing")
		return
	}

	hash := l.Hash(filename)
	l.AllActive[filename] = hash

	l.ensureWorkDirectory(l.WorkDirectory)
	hashFilename := l.hashFilename(hash)

	err := l.AppFs.Rename(filename, hashFilename)
	if err != nil {
		log.WithFields(log.Fields{
			"workDir":  l.WorkDirectory,
			"filename": filename},
		).Fatal("Major problem creating lock file")
	}
}

func (l *locker) Unlock(filename string) {
	hash := l.Hash(filename)
	hashFilename := l.hashFilename(hash)

	err := l.AppFs.Remove(hashFilename)
	if err != nil {
		log.WithFields(log.Fields{
			"workDir":  l.WorkDirectory,
			"filename": filename},
		).Fatal("Major problem unlocking file")
	}
}

func (l *locker) ensureWorkDirectory(path string) {
	exists, err := afero.DirExists(l.AppFs, path)
	if !exists && err == nil {
		l.AppFs.MkdirAll(path, 0755)
	}
}

func (l *locker) hashFilename(hash string) string {
	return filepath.Join(l.WorkDirectory, hash)
}
