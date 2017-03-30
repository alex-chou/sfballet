package main

import (
	"net/http"

	"golang.org/x/net/html"
)

type Request func() (*http.Response, error)

func getRequest(url string) Request {
	return func() (*http.Response, error) {
		return http.Get(url)
	}
}

func requestToNode(request Request) (*html.Node, error) {
	response, err := request()
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return html.Parse(response.Body)
}

