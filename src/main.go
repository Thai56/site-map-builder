/* @Test url for testing the GET requests
 *TEST_URLS: "http://example.com/" http://calhoun.io
 */
package main

import (
	"bytes"
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/gophercises/html-link-parser-4/src"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func isValidUrl(potentialUrl string) bool {
	_, err := url.ParseRequestURI(potentialUrl)
	if err != nil {
		return false
	} else {
		return true
	}
}

func getDomainFromUrl(u string) string {
	domainSlice := strings.Split(u, ".")
	domain := domainSlice[len(domainSlice)-2] + "." + domainSlice[len(domainSlice)-1]
	return domain
}

var visited = make(map[string]bool)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Please specify start page")
		os.Exit(1)
	}

	fmt.Println("args", args)
	queue := make(chan string)

	go func() { queue <- args[0] }()

	for uri := range queue {
		enqueue(uri, queue)
	}
	fmt.Println("REACHING_THE_END")
}

func enqueue(uri string, queue chan string) {
	//TODO: add uri to visited since we just got our return
	visited[uri] = true
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client := http.Client{Transport: transport}

	resp, err := client.Get(uri)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	doc, err := html.Parse(bytes.NewReader(body))

	links, err := linkParser.GetLinks(doc)
	if err != nil {
		fmt.Println("Could not get links", err)
	}

	listOfLinks := linkParser.FormatLinksStruct(links)

	for _, link := range listOfLinks {
		absolute := fixUrl(link.Href, uri)
		if uri == "" {
			fmt.Println("__URI_EMPTY")
		}
		if uri != "" {
			if !visited[absolute] {
				go func() { queue <- absolute }()
			}
		}
	}

	//TODO: check if the visiting address is related to the original argument passed in
	fmt.Println("__VISITED", visited)
}

func fixUrl(href, base string) string {
	uri, err := url.Parse(href)
	if err != nil {
		return ""
	}

	baseUrl, err := url.Parse(base)
	if err != nil {
		return ""
	}
	uri = baseUrl.ResolveReference(uri)
	return uri.String()
}
