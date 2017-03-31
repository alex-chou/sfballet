package main

import (
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
