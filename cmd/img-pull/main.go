// img-pull
// Copyright 2018 Ryan Rasti
// All rights reserved

package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

func main() {
	imagePull(".gif", `C:\img-pull\`, "https://cultofthepartyparrot.com/") // test it out
}

func imagePull(dataType string, downloadDir string, domains ...string) error {
	var err error

	// basic collector
	c := colly.NewCollector(
		// setup domains to pull from
		// colly.AllowedDomains(domains...),

		// set bogus user-agent
		colly.UserAgent("xyz"),

		colly.AllowURLRevisit(),

		// execute collector requests asynchronously
		colly.Async(true),
	)
	d := c.Clone() // duplicate collector to be used for downloading images

	// logs for each request
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Collecting all", dataType, "from", r.URL.String())
	})

	if dataType == ".gif" {
		c.OnHTML("img[src]", func(e *colly.HTMLElement) {
			src := e.Attr("src") // get the src attribute of the img
			fmt.Printf("Image found with src: %s\n", src)

			// handle image src with content hosted at FQDNs
			if strings.Contains(src, "http") {
				e.Request.Visit(src)
				c.OnResponse(func(r *colly.Response) {
					fmt.Printf("Navigated to direct image source!")
				})
			} else {
				fmt.Println("Visiting absolute URI for image src", src, "at", e.Request.AbsoluteURL(src))
				d.Visit(e.Request.AbsoluteURL(src))
				d.OnResponse(func(r *colly.Response) {
					downloadFileFromResponse(e.Request.AbsoluteURL(src), downloadDir, r)
				})
			}
		})
	}

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

	// TODO: understand the async ordering better... the calls below execute in the wrong order
	d.Wait() // wait until download aggregater collector threads are finished
	c.Wait() // wait until main collector thread is finished

	return nil
}

func downloadFileFromResponse(url string, dir string, r *colly.Response) error {
	var err error

	fmt.Println("Saving file at", dir+r.FileName())
	err = r.Save(dir + r.FileName()) // FileName returns the sanitized file name parsed from "Content-Disposition" header or from URL
	if err != nil {
		log.Print("unable to save image from URL", url, "to filesystem location", dir)
		return err
	}

	return nil
}
