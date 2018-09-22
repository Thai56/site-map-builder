package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/gophercises/html-link-parser-4/src"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"net/url"
)

func isValidUrl(potentialUrl string) bool {
	_, err := url.ParseRequestURI(potentialUrl)
	if err != nil {
		return false
	} else {
		return true
	}
}

// TEST_URLS: "http://example.com/" http://calhoun.io
func main() {
	//TODO: refactor main method
	urlFlag := flag.String("url", "http://calhoun.io", "flag for url")

	flag.Parse()

	if urlNotValid := !isValidUrl(*urlFlag); urlNotValid {
		panic(
			fmt.Sprintf(
				"not a valid url please follow the format: http://{your_url_name} | URL%v:",
				*urlFlag,
			),
		)
	}

	res, err := http.Get(*urlFlag)
	if err != nil {
		fmt.Printf("%v\n: ", err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	doc, err := html.Parse(bytes.NewReader(body))

	tester, _ := links.GetLinks(doc)

	fmt.Println("test", links.FormatLinksStruct(tester))
	//TODO: Follow up w/ gopherexercises to get the next instructions
}
