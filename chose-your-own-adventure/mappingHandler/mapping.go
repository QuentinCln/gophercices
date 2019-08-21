package mappingHandler

import (
	model "../parser"
	"net/http"
	"html/template"
)

func NewAdventureHandler(story model.Story, templ *template.Template) http.Handler {
	hand := AdventureHandler{story, templ}
	return hand
}

func (handler AdventureHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if path == "/" {
		handler.HtmlTemplate.Execute(w, handler.Story["intro"])
		return 
	}
	path = path[len("/"):]
	err := handler.HtmlTemplate.Execute(w, handler.Story[path])
	if err != nil{
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
	return
}

type AdventureHandler struct {
	Story model.Story
	HtmlTemplate *template.Template
}
