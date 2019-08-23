package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"./htmllinkparser"
	"strings"
	"encoding/json"
	"io/ioutil"
)

type page struct {
	From string `json:"from"`
	NumberLink int `json:"numberLink"`
	Urls []string `json:"urls"`
}

func removeDuplicate(s []string) []string {
	seen := make(map[string]struct{}, len(s))
	j := 0
	for _, v := range s {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		s[j] = v
		j++
	}
	return s[:j]
}

func pruneUrl(links []htmllinkparser.Link, baseUrl string) []string {
	var keeped, pruned []string
	for _, link := range links {
		if !strings.ContainsRune(link.Href, '#') && strings.HasPrefix(link.Href, "/") {
			link.Href = baseUrl + link.Href
		}
		if strings.HasPrefix(link.Href,baseUrl) {
			keeped = append(keeped, link.Href)
			} else {
				pruned = append(pruned, link.Href)
			}
	}
	fmt.Println("keeped", keeped, "\npruned", pruned)
	return removeDuplicate(keeped)
}

func baseUrl(urlPtr *url.URL) string {
	url := url.URL{
		Scheme: urlPtr.Scheme,
		Host: urlPtr.Host,
	}
	return url.String()
}

func get(url string) ([]htmllinkparser.Link, string) {
	fmt.Println("Mapping :", url)
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	links := htmllinkparser.Parse(response.Body)
	return links, baseUrl(response.Request.URL)
}

func mapSite(basUrl string, baseDepth int) []string {
	var mappedUrls []string
	mappedUrls = append(mappedUrls, basUrl)
	urlSeen := make(map[string]struct{})

	depth:= baseDepth
	for i := 0; i != len(mappedUrls) && depth >= 0; i, depth = i + 1, depth - 1 {
		for _, l := range pruneUrl(get(mappedUrls[i])) {
			if _, ok := urlSeen[l]; !ok {
				urlSeen[l] = struct{}{}
				mappedUrls = append(mappedUrls, l)
			}
		}
	}
	return mappedUrls
}

func writeToJson(siteUrl string, urls []string) {
	ur, err := url.Parse(siteUrl)
	if err != nil {
		panic(err)
	}
	baseUrl := baseUrl(ur)
	page := page{
		From: baseUrl,
		NumberLink: len(urls),
		Urls: urls,
	}
	byteArr, err := json.MarshalIndent(page, " ", " ")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(ur.Host + "-site-map.json", byteArr, 0644)
	if err != nil {
		panic(err)
	}
}

func main() {
	url := flag.String("url", "https://github.com", "URL to parse")
	depth := flag.Int("depth", 40, "DEPTH of sub url")
	flag.Parse()
	urls := mapSite(*url, *depth)
	urls = removeDuplicate(urls) // with double baseUrl getting there baseUrl/ & baseUrl
	fmt.Printf("Found %d urls:\n=> %+v\n", len(urls), urls)
	writeToJson(*url, urls)
}
