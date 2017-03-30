package main

import (
	"strings"

	"golang.org/x/net/html/atom"
	"golang.org/x/net/html"
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

func nodeToProgram(node *html.Node) *Program {
	program := &Program{}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		switch child.DataAtom {
		case atom.A:
			if extractValue(child, "class") == "naked" {
				titleNode := fetch(child, func(node *html.Node) bool {
					return node.DataAtom == atom.Div && extractValue(node, "class") == "large_sub_title"
				})
				program.InfoPath = extractValue(child, "href")
				program.Title = fetchText(titleNode)
			}
		case atom.Ul:
			if extractValue(child, "class") == "performance-list" {
				items := fetchAll(child, byDataAtom(atom.Li))
				program.PerformanceList = make([]string, len(items))
				for idx, item := range items {
					program.PerformanceList[idx] = fetchText(item)
				}
			}
		case atom.P:
			if extractValue(child, "class") == "performance-date" {
				program.Dates = fetchText(child)
			}
		case atom.Div:
			paths := fetchTicketPaths(child)
			if len(paths) > 0 {
				program.TicketPath = paths[0]
			}
		}
	}
	return program
}
