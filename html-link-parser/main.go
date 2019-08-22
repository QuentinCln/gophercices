package main

import (
	"fmt"
	"net/http"
	"golang.org/x/net/html"
	"strings"
	"encoding/json"
	"io/ioutil"
)

// TEST =>
var exampleHtml = `
<html>
<body>
  <h1>Hello!</h1>
  <div>
	<a href="/other-page">
		A link to another page
		<span> some span  </span>
	</a>
  </div>
  <a href="/page-two">A link to a second page
	<div>
	  <a>
		  test
	  </a>
	</div>	
  </a>
  <a></a>
</body>
</html>` 

type Link struct {
	Href string
	Text string
}

func getHtml(url string) *html.Node {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	htmlNode, err := html.Parse(resp.Body)
	if err != nil { // Invalid html
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

func getTextFromNode(node *html.Node) (string) {
	if node.Type == html.TextNode {
		return strings.TrimSpace(node.Data)
	}
	if node.Type != html.ElementNode {
		return "EMPTY"
	}
	var txt string
	for child := node.FirstChild; child != nil; child = child.NextSibling { // go through hitNextLink
		txt += getTextFromNode(child)
	}
	return txt
}

func getHrefFromNode(node *html.Node) string {
	for _, attr := range node.Attr {
		if attr.Key == "href" {
			return attr.Val
		}
	}
	return "EMPTY HREF"
}

func linkNodesToLinks(linkNodes []*html.Node) (links []Link) {
	for _, node := range linkNodes {
		var link Link
		link.Href = getHrefFromNode(node)
		link.Text = getTextFromNode(node)
		links = append(links, link)
	}
	return links
}

func hitNextLink(node *html.Node) bool { // WIP
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if isLink(node) {
			return false
		}
		isLastLink(child)
	}
	return true
}

func isLastLink(node *html.Node) bool { // WIP
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if isLink(node) {
			return false
		}
		isLastLink(child)
	}
	return true
}

func getLinkNodes(document *html.Node, out *int) (nodes []*html.Node) {
	if isLink(document) {
		*out += 1
		return []*html.Node{document}
	}
	for child := document.FirstChild; child != nil; child = child.NextSibling {
		nod := getLinkNodes(child, out) // go trhough isLastLink() before appenning
		for _, n := range nod {
			nodes = append(nodes, n)
		}
	}
	return nodes
}

func writeToJson(links []Link) {
	byteArr, err := json.MarshalIndent(links, "", "")
	if err != nil {
		panic(err)
	}
	er := ioutil.WriteFile("link.json", byteArr, 0644)
	if er != nil {
		panic(err)
	}
}

func main() {
	test := 0

	url := "https://www.google.com/search?q=t&oq=t+&aqs=chrome..69i57j69i59j69i60l2j69i61j69i60.695j0j4&sourceid=chrome&ie=UTF-8"
	links := getLinkNodes(getHtml(url), &test)
	fmt.Println(test)

	// r := strings.NewReader(exampleHtml)
	// htmlNode, err := html.Parse(r)
	// if err != nil {
	// 	panic(err)
	// }
	// links := getLinkNodes(htmlNode, &test)
	
	fmt.Println("Number of links found:", test)
	fmt.Printf("%+v\n", links)
	lks := linkNodesToLinks(links)
	fmt.Printf("%+v\n", lks)
	writeToJson(lks)
}