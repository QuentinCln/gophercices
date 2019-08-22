package htmllinkparser

import (
	"golang.org/x/net/html"
	"strings"
	"encoding/json"
	"io/ioutil"
	"io"
)

type Link struct {
	Href string
	Text string
}

func getTextFromNode(node *html.Node) (string) {
	if node.Type == html.TextNode {
		return strings.TrimSpace(node.Data) + " "
	}
	if node.Type != html.ElementNode {
		return "EMPTY"
	}
	var txt string
	for child := node.FirstChild; child != nil; child = child.NextSibling {
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
	return "EMPTY"
}

func linkNodesToLinks(linkNodes []*html.Node) (links []Link) {
	for _, node := range linkNodes {
		var link Link
		link.Href = getHrefFromNode(node)
		txt := getTextFromNode(node)
		if (len(txt) == 0) {
			link.Text = txt
		} else {
			link.Text = txt[0:len(txt) -1]
		}
		links = append(links, link)
	}
	return links
}


func isLink(doc *html.Node) bool {
	if doc.Type == html.ElementNode && doc.Data == "a" {
		return true
	} 
	return false
}

func getLinkNodes(document *html.Node) (nodes []*html.Node) {
	if isLink(document) {
		return []*html.Node{document}
	}
	for child := document.FirstChild; child != nil; child = child.NextSibling {
		nod := getLinkNodes(child)
		for _, n := range nod {
			nodes = append(nodes, n)
		}
	}
	return nodes
}

func WriteToJson(links []Link) {
	byteArr, err := json.MarshalIndent(links, "", "")
	if err != nil {
		panic(err)
	}
	er := ioutil.WriteFile("link.json", byteArr, 0644)
	if er != nil {
		panic(err)
	}
}

func Parse(r io.Reader) []Link {
	htmlNode, err := html.Parse(r)
	if err != nil {
		panic(err)
	}
	nodes := getLinkNodes(htmlNode)
	links := linkNodesToLinks(nodes)
	return links
}