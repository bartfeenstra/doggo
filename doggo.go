package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
	"time"
)

var impatientClient = &http.Client{
	// Enforce a timeout so remotes cannot block our own process.
	Timeout: time.Second * 10,
}

type ImageData struct {
	OriginalUrl string `json:"image_original_url"`
}

type Image struct {
	Data ImageData
}

func respond(writer http.ResponseWriter, request *http.Request) {
	// Retrieve the metadata for a random image.
	response, err := impatientClient.Get("https://api.giphy.com/v1/gifs/random?tag=puppy&api_key=dc6zaTOxFJmzC")
	if err != nil {
		log.Fatal(err)
		return
	}

	// Extract the image's URL.
	responseJson, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		return
	}
	image := Image{}
	responseJsonBytes := []byte(responseJson)
	err = json.Unmarshal(responseJsonBytes, &image)
	if err != nil {
		log.Fatal("Invalid Giphy response JSON: ", err)
		return
	}
	pageTemplate := template.Must(template.ParseFiles("templates/page.tmpl"))
	pageTemplate.Execute(writer, image)
}

func main() {
	http.HandleFunc("/", respond)
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal("HTTP server error: ", err)
	}
}
