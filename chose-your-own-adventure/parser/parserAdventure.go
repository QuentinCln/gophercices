package parser

import (
	"encoding/json"
	"io/ioutil"
)

type Story map[string]Adventure

type Adventure struct {
	Title string `json:"title"`
	Descriptions []string `json:"story"` 
	Options []Option`json:"options"`		
}

type Option struct {
	Description string `json:"text"`
	Next string `json:"arc"`
}


func ParseAdventure(filename string) (story Story, err error) {
	jsonData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonData, &story)
	if err != nil {
		return nil, err
	}
	return story, nil
}