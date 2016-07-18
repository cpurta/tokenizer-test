package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

func main() {
	url := "https://www.youtube.com/watch?v=0BPUVxuAYf4"
	// grab the response from the site...
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("Error getting response from", url, "Error:", err.Error())
	}

	if resp.Body != nil {
		defer resp.Body.Close()

		tokenizer := html.NewTokenizer(resp.Body)

		startTime := time.Now()

		depth := 0
		count := 1
		for {
			tt := tokenizer.Next()
			switch tt {
			case html.ErrorToken:
				// should only happen when we hit EOF
				goto finished
			case html.TextToken:
				if depth > 0 {
					fields := strings.Fields(string(tokenizer.Text()))
					for _, f := range fields {
						if f == "views" {
							fmt.Printf("Count %d: %q\n\n", count, fields)
							count++
						}
					}
				}
			case html.StartTagToken, html.EndTagToken:
				tn, _ := tokenizer.TagName()
				if len(tn) == 3 && string(tn[:3]) == "div" {
					if tt == html.StartTagToken {
						depth++
					} else {
						depth--
					}
				}
			} // end switch
		} // end for

	finished:
		fmt.Printf("Finished parsing %s in %.2f seconds\n", url, time.Since(startTime).Seconds())
	} // end if
}
