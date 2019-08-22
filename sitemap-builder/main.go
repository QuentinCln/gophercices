package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"./htmllinkparser"
	"strings"
)

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

func mapSite(basUrl string) []string {
	var mappedUrls []string
	mappedUrls = append(mappedUrls, basUrl)
	urlSeen := make(map[string]struct{})
	
	for i := 0; i != len(mappedUrls); i++ {
		for _, l := range pruneUrl(get(mappedUrls[i])) {
			if _, ok := urlSeen[l]; !ok {
				urlSeen[l] = struct{}{}
				mappedUrls = append(mappedUrls, l)
			}
		}
	}
	return mappedUrls
}

func main() {
	url := flag.String("url", "https://discordapp.com/", "URL to parse")
	flag.Parse()
	urls := mapSite(*url)
	urls = removeDuplicate(urls) // to be sure
	fmt.Printf("Found %d mapped links / %+v", len(urls), urls)
}
