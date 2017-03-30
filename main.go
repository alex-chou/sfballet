package main

import (
	"fmt"
	"log"
	"strings"

	"golang.org/x/net/html"
)

const (
	sfballetRoot = "https://www.sfballet.org"
)

func main() {
	fmt.Println("fetching ticket paths")
	doc, err := requestToNode(getRequest(sfballetRoot))
	if err != nil {
		log.Fatal(err)
	}
	paths := fetchTicketPaths(doc)

	fmt.Println("fetching programs")
	for _, path := range paths {
		fmt.Println(path)
		doc, err = requestToNode(getRequest(fmt.Sprintf("%s%s", sfballetRoot, path)))
		if err != nil {
			log.Fatal(err)
		}
		programs := fetchPrograms(doc)
		for _, program := range programs {
			fmt.Printf("Program: %s\n", program.Title)
			fmt.Printf("Info Link: %s%s\n", sfballetRoot, program.InfoPath)
			if len(program.PerformanceList) > 0 {
				fmt.Printf("Performance List: %s\n", strings.Join(program.PerformanceList, ", "))
			}
			if program.Dates == "" {
				fmt.Println("Dates: This program is no longer available.")
				fmt.Println("Buy Tickets Link: This program is no longer available.")
			} else {
				fmt.Printf("Dates: %s\n", program.Dates)
				fmt.Printf("Buy Tickets Link: %s%s\n", sfballetRoot, program.TicketPath)
			}
		}
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
