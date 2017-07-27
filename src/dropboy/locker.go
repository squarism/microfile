package dropboy

import (
	"crypto/md5"
	"fmt"
	"io"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type locker struct {
	AllActive     map[string]string
	AppFs         afero.Fs
	WorkDirectory string
}

func NewLocker(workDir string) locker {
	l := locker{}
	l.AllActive = make(map[string]string)
	l.AppFs = afero.NewOsFs()
	l.WorkDirectory = workDir
	return l
}

func (l *locker) Hash(filename string) string {
	hasher := md5.New()
	io.WriteString(hasher, filename)
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

func (l *locker) Lock(filename string, hash string) {
	if _, ok := l.AllActive[filename]; ok {
		log.WithFields(log.Fields{"workDir": l.WorkDirectory, "filename": filename}).Info("Already processing")
	} else {
		l.AllActive[filename] = hash
	}

	l.ensureWorkDirectory(l.WorkDirectory)

	l.AppFs.Rename(filename, l.hashFilename(hash))
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
