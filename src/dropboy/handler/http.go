package handler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	url "net/url"

	"github.com/fsnotify/fsnotify"

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
	req, err := http.NewRequest("POST", destinationURL, bytes.NewBuffer(json))
	if err != nil {
		log.Fatal("Can't make sense of destinationURL: %s", destinationURL)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Printf("HTTP error posting to: %s\n", destinationURL)
	}

	if response != nil {
		defer response.Body.Close()
	} else {
		log.Printf("Empty HTTP response for: %s\n", destinationURL)
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
		log.Fatal("URL %s doesn't appear to be a URL.", u)
	}

	du, err := url.Parse(defaultURL)
	if err != nil {
		log.Fatal("What is configured as DefaultURL %s doesn't appear to be a URL.", u)
	}

	if u.Host == "" {
		u.Host = du.Host
		u.Scheme = du.Scheme
	}

	return u.String()
}
