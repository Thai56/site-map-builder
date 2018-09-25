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

func getUrlFlag() *string {
	urlFlag := flag.String("url", "http://www.calhoun.io", "flag for url")

	flag.Parse()

	if urlNotValid := !isValidUrl(*urlFlag); urlNotValid {
		panic(
			fmt.Sprintf(
				"not a valid url please follow the format: http://www.{your_url_name} | URL%v:",
				*urlFlag,
			),
		)
	}

	return urlFlag
}

func getDomainFromUrl(u string) string {
	domainSlice := strings.Split(u, ".")
	domain := domainSlice[len(domainSlice)-2] + "." + domainSlice[len(domainSlice)-1]
	return domain
}

//TODO: make urlMap a type struct and add this method to it
func updateUrlMapIfKeyDoesNotExist(m map[string]string, target string, title string) { /*TODO: Take out title if necessary*/
	if _, ok := m[target]; !ok {
		m[target] = title
	}
}

// TEST_URLS: "http://example.com/" http://calhoun.io
func main() {
	urlFlag := getUrlFlag()

	res, err := http.Get(*urlFlag)
	if err != nil {
		fmt.Printf("%v\n: ", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	doc, err := html.Parse(bytes.NewReader(body))

	links, err := linkParser.GetLinks(doc)
	if err != nil {
		fmt.Printf("Could not get links %v:", err)
	}

	listOfLinks := linkParser.FormatLinksStruct(links)

	mainDomain := getDomainFromUrl(*urlFlag)

	urlMap := make(map[string]string)

	for _, l := range listOfLinks {
		if firstChar := string(l.Href[0]); firstChar == "/" {
			updateUrlMapIfKeyDoesNotExist(urlMap, l.Href, l.Title)
		} else if strings.Contains(l.Href, mainDomain) {
			updateUrlMapIfKeyDoesNotExist(urlMap, l.Href, l.Title)
		}
	}

	fmt.Println("urlMap \n", urlMap)

	//TODO: MAIN WORK IS BELOW HERE
	//TODO: loop or recurse through each link in the map
	//TODO: visit the link and create a url list out of the links
	//TODO: start off by just doing the first link in the map
	for k, _ := range urlMap {
		if string(k[0]) == "/" {
			res, err := http.Get("http://www." + mainDomain + k)
			if err != nil {
				fmt.Printf("Could not get the domain of subdomain") // TODO: make this error clearer
			}
			defer res.Body.Close()

			//TODO: refactor this method and the one above
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				fmt.Printf("could not get body %v", err)
			}
			//TODO: Handle the error here
			doc, err := html.Parse(bytes.NewReader(body))
			//TODO: Handle the error here
			links, _ := linkParser.GetLinks(doc)
			fmt.Println("Got tha body!!!", linkParser.FormatLinksStruct(links))
		} /*TODO: else just get the url given because it must contain the domain name already*/
	}
}
