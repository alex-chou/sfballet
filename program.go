package main

import (
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Program struct {
	InfoURL         string
	Title           string
	PerformanceList []string
	Available       bool
	Dates           string
	ShowDates       []*Date
	TicketURL       string
}

func (p *Program) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Program: %s\n", p.Title))
	if !p.Available {
		buffer.WriteString("Note: This program is no longer available.\n")
	}
	buffer.WriteString(fmt.Sprintf("Info Link: %s\n", p.InfoURL))
	if len(p.PerformanceList) > 0 {
		buffer.WriteString(fmt.Sprintf("Performance List: %s\n", strings.Join(p.PerformanceList, ", ")))
	}
	if p.Available {
		buffer.WriteString(fmt.Sprintf("Dates: %s\n", p.Dates))
		buffer.WriteString(fmt.Sprintf("Buy Tickets Link: %s\n", p.TicketURL))
	}
	return buffer.String()
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
				program.InfoURL = fmt.Sprintf("%s%s", sfballetRoot, extractValue(child, "href"))
				program.Title = fetchText(titleNode)
			}
		case atom.Ul:
			if extractValue(child, "class") == "performance-list" {
				items := fetchAll(child, byDataAtom(atom.Li))
				program.PerformanceList = make([]string, len(items))
				for idx, item := range items {
					program.PerformanceList[idx] = strings.TrimSpace(fetchText(item))
				}
			}
		case atom.P:
			if extractValue(child, "class") == "performance-date" {
				program.Available = true
				program.Dates = fetchText(child)
			}
		case atom.Div:
			paths := fetchTicketPaths(child)
			if len(paths) > 0 {
				program.TicketURL = fmt.Sprintf("%s%s", sfballetRoot, paths[0])
			}
		}
	}
	return program
}
