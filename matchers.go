package main

import (
	"strings"

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

func fetchTicketPaths(node *html.Node) []string {
	tMatcher := func(node *html.Node) bool {
		if node.DataAtom == atom.A {
			textNodes := fetchAll(node, textMatcher)
			for _, textNode := range textNodes {
				if strings.Contains( strings.ToLower(textNode.Data), "buy tickets") {
					return true
				}
			}
		}
		return false
	}
	matches := fetchAll(node, tMatcher)
	urls := make([]string, len(matches))
	for idx, match := range matches {
		urls[idx] = extractValue(match, "href")
	}
	return urls
}

func fetchPrograms(node *html.Node) []*Program {
	pMatcher := func(node *html.Node) bool {
		return node.DataAtom == atom.Div && extractValue(node, "class") == "program-item"
	}
	matches := fetchAll(node, pMatcher)
	programs := make([]*Program, len(matches))
	for idx, match := range matches {
		programs[idx] = nodeToProgram(match)
	}
	return programs
}
