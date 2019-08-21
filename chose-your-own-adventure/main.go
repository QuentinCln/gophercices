package main

import (
	"os"
	"flag"
	"fmt"
	"net/http"
	"html/template"
	parserAdventure "./parser"
	"./mappingHandler"
)


func main() {
	fileName := flag.String("filename", "adventure.json", "JSON file that contains the adventure")
	htmlFile := flag.String("htmlfile", "adventure.html", "HTML file that contains the template")
	
	htmlTemplate := template.Must(template.ParseFiles(*htmlFile))
	story, err := parserAdventure.ParseAdventure(*fileName)
	if err != nil {
		fmt.Println("error", err)
		fmt.Println(htmlTemplate)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.Handle("/", mappingHandler.NewAdventureHandler(story, htmlTemplate))
	if err != nil {
		fmt.Println("handler error", err)
		os.Exit(1)
	}

	http.ListenAndServe(":8080", mux)
	os.Exit(0)
}






