// img-pull
// Copyright 2018 Ryan Rasti
// Refer to IDEAS.MD

package main

import (
	"fmt"
	"log"

	"github.com/gocolly/colly"
)

func main() {
	// ultimate test case: https://cultofthepartyparrot.com/
	imagePull(".gif", "https://cultofthepartyparrot.com/")
}

// the workhorse
// leverages the colly framework http://go-colly.org/
// TODO: Remove me later
// https://github.com/gocolly/colly/blob/master/_examples/basic/basic.go
func imagePull(dataType string, domains ...string) error {
	var err error

	c := colly.NewCollector(
		// setup domains to pull from
		colly.AllowedDomains(domains...),
	)

	// callback for each request
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Collecting", dataType, "from", r.URL.String())
	})

	if dataType == ".gif" {
		c.OnHTML("img[src]", func(e *colly.HTMLElement) {
			src := e.Attr("src") // get the src attribute of the img
			fmt.Printf("Image found with src: %s", src)
		})
	}

	for _, v := range domains {

		// collect data from each domain
		err = c.Visit(v)
		if err != nil {
			l := fmt.Sprintf("error visiting %s", v)
			log.Fatal(l, err)
			return err
		}
	}

	return nil
}
