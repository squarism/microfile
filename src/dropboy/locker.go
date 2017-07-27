package dropboy

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"

	"github.com/spf13/afero"
)

type locker struct {
	AllActive map[string]string
	AppFs     afero.Fs
}

func NewLocker() locker {
	l := locker{}
	l.AllActive = make(map[string]string)
	l.AppFs = afero.NewOsFs()
	return l
}

func (l *locker) Hash(filename string) string {
	hasher := md5.New()
	io.WriteString(hasher, filename)
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

func (l *locker) Lock(filename string, hash string) error {
	if _, ok := l.AllActive[filename]; ok {
		return errors.New(fmt.Sprintf("Already processing %s", filename))
	} else {
		l.AllActive[filename] = hash
		return nil
	}
}

func (l *locker) EnsureWorkDirectory(path string) {
	exists, err := afero.DirExists(l.AppFs, path)
	if !exists && err == nil {
		l.AppFs.MkdirAll(path, 0755)
	}
}
