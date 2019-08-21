package main

import (
	"fmt"
	"net/http"
	"./urlshort"
	"io/ioutil"
	"os"
	
)

func readFile(filename string) []byte {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		os.Exit(1)
	}
	return data
}

func main() {
	jsonFile := "url.json"
	yamlFile := "url.yaml"

	mux := defaultMux()

	yamlHandler, err := urlshort.YAMLHandler(readFile(yamlFile), mux)
	if err != nil {
		panic(err)
	}
	
	jsonHandler, err := urlshort.JSONHandler(readFile(jsonFile), yamlHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, world!")
	})
	return mux
}