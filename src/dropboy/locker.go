package dropboy

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
)

type locker struct {
	AllActive map[string]string
}

func NewLocker() locker {
	l := locker{}
	l.AllActive = make(map[string]string)
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
