package handler

import (
	"log"
	"path/filepath"
	"testing"

	"github.com/fsnotify/fsnotify"
	"github.com/stretchr/testify/assert"
	httpmock "gopkg.in/jarcoal/httpmock.v1"

	"dropboy/config"
)

func TestHttpSendEvent(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST",
		"http://localhost:3000/api/dostuff",
		httpmock.NewStringResponder(200, `OK`),
	)
	httpHandler := HTTP{DefaultURL: "http://localhost:3000/", Path: "/api/dostuff", SendContents: false}

	filename, err := filepath.Abs("./fixtures/dropbox/wee.txt")
	if err != nil {
		log.Fatal(err)
	}
	event := fsnotify.Event{Name: filename, Op: fsnotify.Create}
	httpHandler.Handle(event)

	calls := httpmock.GetTotalCallCount()

	assert.Equal(t, 1, calls)
}

func TestUsesDefaultWhenHostIsBlank(t *testing.T) {
	httpHandler := HTTP{DefaultURL: "http://my.server/"}
	url := httpHandler.pathCompletion("/api")

	expected := "http://my.server/api"

	assert.Equal(t, expected, url)
}

func TestUsesDefaultWhenFullURL(t *testing.T) {
	httpHandler := HTTP{DefaultURL: "http://my.server/"}
	url := httpHandler.pathCompletion("http://thingy.biz/api")

	expected := "http://thingy.biz/api"

	assert.Equal(t, expected, url)
}

func TestConfiguresPath(t *testing.T) {
	httpHandler := new(HTTP)
	action := config.Action{
		Type: "http",
		Options: map[string]string{
			"path": "/images",
		},
	}
	httpHandler.Init(action)

	assert.Equal(t, "/images", httpHandler.Path)
}
