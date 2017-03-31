package main

import (
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func fetchTicketPaths(node *html.Node) []string {
	tMatcher := func(node *html.Node) bool {
		if node.DataAtom == atom.A {
			textNodes := fetchAll(node, textMatcher)
			for _, textNode := range textNodes {
				if strings.Contains(strings.ToLower(textNode.Data), "buy tickets") {
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
