package handler

import (
	"io/ioutil"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"

	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"microfile/config"
)

type Imaginary struct {
	Path            string
	DefaultURL      string
	Type            string
	Quality         int
	OutputDirectory string
}

// We can't process images from these events probably
func (i *Imaginary) ignoredEvent(event fsnotify.Event) bool {
	return (event.Op == fsnotify.Remove || event.Op == fsnotify.Chmod)
}

func (i Imaginary) Handle(event fsnotify.Event) {
	if i.ignoredEvent(event) == true {
		return
	}

	file, _ := os.Open(event.Name)
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", filepath.Base(file.Name()))
	io.Copy(part, file)
	writer.Close()

	destinationURL := PathCompletion(i.Path, i.DefaultURL)

	r, _ := http.NewRequest("POST", destinationURL, body)
	r.Header.Add("Accept", writer.FormDataContentType())
	r.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}

	response, err := client.Do(r)
	if err != nil {
		log.WithFields(
			log.Fields{
				"handler":      "imaginary",
				"changed_file": event.Name,
				"url":          destinationURL,
			}).Warn("Problem with imaginary service")
		return
	}

	if response.StatusCode == 200 {
		filename, err := i.writeImage(response, event.Name)
		if err != nil {
			log.WithFields(
				log.Fields{
					"handler":      "imaginary",
					"changed_file": event.Name,
					"output_file":  filename,
				}).Warn("Problem saving converted image")
		}
	} else {
		log.WithFields(
			log.Fields{
				"handler":      "imaginary",
				"status_code":  response.StatusCode,
				"changed_file": event.Name,
			}).Warn("Imaginary response")
	}

}

func (i *Imaginary) Init(action config.Action) error {
	if action.Options["path"] != "" {
		i.Path = action.Options["path"]
	}

	outputDirectory := action.Options["output_directory"]
	if outputDirectory != "" {
		i.OutputDirectory = outputDirectory
	} else {
		log.WithFields(
			log.Fields{
				"handler": "imaginary",
			}).Fatal("output_directory for Imaginary needs to be set (where converted images are saved)")
	}
	return nil
}

func (i *Imaginary) writeImage(response *http.Response, filePath string) (string, error) {
	basename := filepath.Base(filePath)
	destinationFile := filepath.Join(i.OutputDirectory, basename)

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.WithFields(
			log.Fields{
				"handler":     "imaginary",
				"output_file": destinationFile,
			}).Fatal("Can't read response Body")
	}
	defer response.Body.Close()

	err = ioutil.WriteFile(destinationFile, bodyBytes, 0644)
	if err != nil {
		log.WithFields(
			log.Fields{
				"handler":          "imaginary",
				"output_file":      destinationFile,
				"output_directory": i.OutputDirectory,
			}).Warn("Can't write file.  Does output_directory exist?")
	} else {
		log.WithFields(
			log.Fields{
				"handler":     "imaginary",
				"output_file": destinationFile,
			}).Info("Converted image.")
	}

	return destinationFile, nil
}
