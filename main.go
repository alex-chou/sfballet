package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

const (
	sfballetRoot = "https://www.sfballet.org"
)

var hrefKey = map[string]bool{"href": true}

func main() {
	response, err := http.Get(sfballetRoot)
	if err != nil {
		log.Fatal(err)
	}

	hrefs := fetchAll(response, "buy tickets", hrefKey)
	fmt.Println(hrefs)
}

func fetchAll(response *http.Response, text string, keys map[string]bool) []string {
	results := []string{}
	tokenizer := html.NewTokenizer(response.Body)
	defer func() {
		response.Body.Close()
	}()

	for {
		switch tNext := tokenizer.Next(); tNext {
		case html.ErrorToken:
			return results
		case html.StartTagToken:
			tStart := tokenizer.Token()
			if isAnchor := tStart.Data == "a"; isAnchor {
				switch tNext := tokenizer.Next(); tNext {
				case html.TextToken:
					tText := tokenizer.Token()
					if strings.Contains(
						strings.ToLower(tText.String()),
						strings.ToLower(text),
					) {
						results = append(results, extractAttributeValues(tStart, keys)...)
					}
				}
			}
		}
	}
	return results
}

func extractAttributeValues(token html.Token, keys map[string]bool) []string {
	if len(keys) == 0 {
		return nil
	}
	values := []string{}
	for _, attr := range token.Attr {
		if _, ok := keys[attr.Key]; ok {
			values = append(values, attr.Val)
		}
	}
	return values
}
