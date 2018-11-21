// img-pull
// Copyright 2018 Ryan Rasti
// All rights reserved

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gocolly/colly"
)

func main() {
	imagePull(".gif", `C:\img-pull\`, "https://cultofthepartyparrot.com/") // test it out
}

func imagePull(dataType string, downloadDir string, domains ...string) error {
	var err error

	c := colly.NewCollector(

		// uncomment to enable debugger
		// colly.Debugger(&debug.LogDebugger{}),

		// setup domains to pull from
		// colly.AllowedDomains(domains...),

		// set bogus user-agent
		colly.UserAgent("xyz"),

		colly.AllowURLRevisit(),

		// execute collector requests asynchronously
		colly.Async(true),
	)

	// map to hold absolute URLs of images
	m := map[string]bool{}

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Collecting all", dataType, "from", r.URL.String())
	})

	c.OnHTML("img[src]", func(e *colly.HTMLElement) {
		src := e.Attr("src")
		fmt.Printf("Image found with src: %s\n", src)

		// handle image src with content hosted at FQDNs
		if strings.Contains(src, "http") {
			e.Request.Visit(src)
			c.OnResponse(func(r *colly.Response) {
				fmt.Printf("Navigated to direct image source!")
			})
		} else if strings.Contains(src, ".gif"){
			fmt.Println("Appending absolute URL for image src to array", src, "at", e.Request.AbsoluteURL(src))

			// if the image absolute URL is not contained within the map, add it
			if _, ok := m[e.Request.AbsoluteURL(src)]; !ok {
				m[e.Request.AbsoluteURL(src)] = true
				fmt.Println(e.Request.AbsoluteURL(src), "was added to the map")
			}
		}
	})

	// visit all domains passed into imagePull()
	for _, v := range domains {
		fmt.Println("Visiting", v)

		// collect data from each domain
		err = c.Visit(v)
		if err != nil {
			l := fmt.Sprintf("error visiting %s ", v)
			log.Fatal(l, err)
			return err
		}
	}

	c.Wait() // wait until main collector thread is finished

	// download data using net/http to avoid dealing with nested callbacks in colly
	downloadImageDataFromMap(m, downloadDir)

	return nil
}

func downloadImageDataFromMap(m map[string]bool, dir string) error {
	// traverse the map, and use path.Base to get the "last element of the path" i.e. the filename
	for k, _ := range m {
		fmt.Println("Saving", path.Base(k), "to", dir)

		// create the file
		out, err := os.Create(dir + path.Base(k))
		if err != nil {
			return err
		}
		defer out.Close()

		// get the raw HTTP response data
		res, err := http.Get(k)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		// write the response body (image data) to file
		_, err = io.Copy(out, res.Body)
		if err != nil {
			return err
		}
	}

	return nil
}