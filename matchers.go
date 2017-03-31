package main

import (
	"fmt"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Matcher func(node *html.Node) bool

func textMatcher(node *html.Node) bool {
	return node.Type == html.TextNode
}

func byDataAtom(a atom.Atom) Matcher {
	return func(node *html.Node) bool {
		return node.DataAtom == a
	}
}

func fetch(node *html.Node, matcher Matcher) *html.Node {
	if matcher(node) {
		return node
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if match := fetch(child, matcher); match != nil {
			return match
		}
	}
	return nil
}

func fetchText(node *html.Node) string {
	var text string
	textNodes := fetchAll(node, textMatcher)
	for _, textNode := range textNodes {
		text = fmt.Sprintf("%s%s", text, textNode.Data)
	}
	return text
}

func fetchAll(node *html.Node, matcher Matcher) []*html.Node {
	if matcher(node) {
		return []*html.Node{node}
	}

	results := []*html.Node{}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		results = append(results, fetchAll(child, matcher)...)
	}
	return results
}

func extractValue(node *html.Node, key string) string {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}
