package main

import (
	"fmt"
	"net/http"
	"golang.org/x/net/html"
	"strings"
)

// TEST =>
var exampleHtml = `
<html>
<body>
  <h1>Hello!</h1>
  <a href="/other-page">
    A link to another page
    <span> some span  </span>
  </a>
  <a href="/page-two">A link to a second page
	  <a>
	  	test
	  </a>
  </a>
</body>
</html>` 

type Link struct {
	Href string
	Text string
	Next *Link // for sublink
}

func getHtml(url string) *html.Node {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	htmlNode, err := html.Parse(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(*htmlNode)
	return htmlNode
}

func isLink(doc *html.Node) bool {
	if doc.Type == html.ElementNode && doc.Data == "a" {
		return true
	} 
	return false
}

func getTextFromNode(doc *html.Node) (string, bool) {
	if doc.Type == html.TextNode {
		fmt.Println("###########################")
		return doc.Data, true
	}
	return "", false
}

func nodeToSLink(document *html.Node, out *int, outLinks *[]Link) {
	if isLink(document) {
		fmt.Println("Got <a>")
		if data, ok := getTextFromNode(document); ok {
			fmt.Println("EEZAEAEZAEZA")
			*outLinks = append(*outLinks, Link{"<a>", data, nil})
		}
		*out += 1
	}
	for child := document.FirstChild; child != nil; child = child.NextSibling {
		fmt.Println("child:", child)
		nodeToSLink(child, out, outLinks)
	}
}

func main() {
	// link1 := Link{"<a></a>", "to toto", nil}
	// link2 := Link{"<a></a>", "to toto", &link1}

	// var links []Link
	// links = append(links, link1, link2)
	// fmt.Println(link1)
	// fmt.Println(link2)
	// fmt.Println(links)


	test := 0
	var links []Link
	// url := "https://www.google.com/search?q=t&oq=t+&aqs=chrome..69i57j69i59j69i60l2j69i61j69i60.695j0j4&sourceid=chrome&ie=UTF-8"
	// nodeToSLink(getHtml(url), &test)
	// fmt.Println(test)


	r := strings.NewReader(exampleHtml)
	htmlNode, err := html.Parse(r)
	if err != nil {
		panic(err)
	}

	nodeToSLink(htmlNode, &test, &links)
	fmt.Println(test)


	fmt.Printf("%+v\n", links)
}