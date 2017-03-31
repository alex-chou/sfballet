package main

import (
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Date struct {
	ID      string
	Date    string
	Day     string
	Time    string
	Comment string
}

func (d *Date) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Date: %s %s @ %s\n", d.Day, d.Date, d.Time))
	if d.Comment != "" {
		buffer.WriteString(fmt.Sprintf("\tAdditional comments: %s\n", d.Comment))
	}
	return buffer.String()
}

func fetchShowDates(node *html.Node) []*Date {
	iMatcher := func(node *html.Node) bool {
		return node.DataAtom == atom.Li && strings.Contains(extractValue(node, "class"), "performance-list-item")
	}
	matches := fetchAll(node, iMatcher)
	dates := make([]*Date, len(matches))
	for idx, match := range matches {
		dates[idx] = nodeToDate(match)
	}
	return dates
}

func nodeToDate(parent *html.Node) *Date {
	date := &Date{
		ID:  extractValue(parent, "data-perf-no"),
		Day: strings.Title(extractValue(parent, "data-day")),
	}
	date.Date = strings.TrimSpace(fetchText(fetch(parent, func(node *html.Node) bool {
		return node.DataAtom == atom.H3 && strings.Contains(extractValue(node, "class"), "performance-item-date")
	})))
	date.Time = strings.TrimSpace(fetchText(fetch(parent, func(node *html.Node) bool {
		return node.DataAtom == atom.Div && strings.Contains(extractValue(node, "class"), "performance-item-tod")
	})))
	date.Comment = strings.TrimSpace(fetchText(fetch(parent, func(node *html.Node) bool {
		return node.DataAtom == atom.Div && strings.Contains(extractValue(node, "class"), "performance-item-comment")
	})))
	return date
}
