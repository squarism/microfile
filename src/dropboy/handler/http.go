package handler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	url "net/url"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"

	"dropboy/config"
)

type HTTP struct {
	Path         string
	DefaultURL   string
	SendContents bool
}

type payload struct {
	Filename string `json:"filename"`
	Event    string `json:"event"`
	Contents []byte `json:"contents,omitempty"`
}

func (h HTTP) Handle(event fsnotify.Event) {
	postPayload := payload{
		Filename: event.Name,
		Event:    event.Op.String(),
	}
	if h.SendContents == true {
		contents, _ := ioutil.ReadFile(event.Name)
		postPayload.Contents = contents
	}
	json, _ := json.Marshal(postPayload)

	destinationURL := PathCompletion(h.Path, h.DefaultURL)

	// at this point, the request has not been fired, so err here would be URL parsing errors
	req, _ := http.NewRequest("POST", destinationURL, bytes.NewBuffer(json))

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.WithFields(
			log.Fields{
				"handler": "http",
				"url":     destinationURL,
			}).Warn("Problem while doing an http POST")

		return
	}

	if response != nil {
		defer response.Body.Close()
	} else {
		log.WithFields(
			log.Fields{
				"handler": "http",
				"url":     destinationURL,
			}).Warn("Empty HTTP response")
	}
}

func (h *HTTP) Init(action config.Action) error {
	if action.Options["path"] != "" {
		h.Path = action.Options["path"]
	}
	return nil
}

// pathCompletion completes a URL is needed
// if url is a partial path like /api then the path should
// use the DefaultURL string from config
func PathCompletion(s string, defaultURL string) string {
	u, err := url.Parse(s)
	// TODO: this needs to go into the config validation
	if err != nil {
		log.WithFields(
			log.Fields{
				"handler": "http",
				"url":     u,
			}).Warn("URL doesn't appear to be a URL.")
	}

	du, err := url.Parse(defaultURL)
	if err != nil {
		log.WithFields(
			log.Fields{
				"handler": "http",
				"url":     u,
			}).Warn("What is configured as DefaultURL doesn't appear to be a URL.")
	}

	if u.Host == "" {
		u.Host = du.Host
		u.Scheme = du.Scheme
	}

	return u.String()
}
