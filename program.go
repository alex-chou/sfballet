package main

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Program struct {
	InfoPath string
	Title string
	PerformanceList []string
	Dates string
	TicketPath string
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
