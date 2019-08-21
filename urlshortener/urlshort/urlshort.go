package urlshort

import (
	"net/http"
	yaml "gopkg.in/yaml.v2"
	"encoding/json"
)

func defaultHandler(getMap func(data[]byte) (map[string]string, error), data []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathToUrls, err := getMap(data)
	if err != nil {
		return nil, err
	}
	handler := func(w http.ResponseWriter, r *http.Request) {
		dest, urlFound := pathToUrls[r.URL.Path]
		if urlFound {
			http.Redirect(w, r, dest, http.StatusFound)
		}
		fallback.ServeHTTP(w, r)
	}
	return handler, nil
}

type urlshortener struct {
	// Struct fields are only unmarshalled if they are exported (have an upper case first letter)
	Path string `yaml:"path"`
	Url string `yaml:"url"`
}

func getMapFromYaml(data []byte) (map[string]string, error) {
	var urlMap []urlshortener
	err := yaml.Unmarshal(data, &urlMap)
	if err != nil {
		return nil, err
	}
	pathToUrlMap := make(map[string]string, len(urlMap))
	for _, line := range urlMap {
		pathToUrlMap[line.Path] = line.Url
	}
	return pathToUrlMap, nil
}

func getMapFromJson(data []byte) (map[string]string, error) {
	var urlMap []urlshortener
	err := json.Unmarshal(data, &urlMap)
	if err != nil {
		return nil, err
	}
	pathToUrlMap := make(map[string]string, len(urlMap))
	for _, line := range urlMap {
		pathToUrlMap[line.Path] = line.Url
	}
	return pathToUrlMap, nil
}


func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	return defaultHandler(getMapFromYaml, yml, fallback)
}

func JSONHandler(json []byte, fallback http.Handler) (http.HandlerFunc, error) {
	return defaultHandler(getMapFromJson, json, fallback)
}