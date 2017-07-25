package handler

import (
	"io/ioutil"
	"log"

	"github.com/fsnotify/fsnotify"

	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"dropboy/config"
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
		log.Printf("Problem with imaginary service: %s\n", err)
		os.Exit(1)
	}

	if response.StatusCode == 200 {
		_, err = i.writeImage(response, event.Name)
		if err != nil {
			log.Fatal("Problem saving converted image: %s", err)
		}
	} else {
		log.Println("Imaginary couldn't convert the image %s", file.Name())
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
		log.Fatal("output_directory for Imaginary needs to be set (where converted images are saved).")
	}
	return nil
}

func (i *Imaginary) writeImage(response *http.Response, filePath string) (string, error) {
	basename := filepath.Base(filePath)
	destinationFile := filepath.Join(i.OutputDirectory, basename)

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Can't read response Body wtf: %s", err)
	}
	defer response.Body.Close()

	err = ioutil.WriteFile(destinationFile, bodyBytes, 0644)

	return destinationFile, nil
}
