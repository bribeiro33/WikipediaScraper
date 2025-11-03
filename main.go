package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

// each Page will become one JSON object in the NDJSON file
type Page struct {
	URL       string    `json:"url"`
	Text      string    `json:"text"`
	CrawledAt time.Time `json:"crawled_at"`
}

func main() {
	// expect: go run main.go <input_file> <output_file>
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <urls.txt> <output.ndjson>")
		os.Exit(1)
	}
	inputFile := os.Args[1]
	outputFile := os.Args[2]

	// read urls from txt file
	urls, err := readURLs(inputFile)
	if err != nil {
		log.Fatalf("Error reading URL file: %v", err)
	}

	// create the output file
	file, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}
	defer file.Close()

	// create a Colly collector
	c := colly.NewCollector(
		colly.AllowedDomains("en.wikipedia.org"), // restrict crawler to Wikipedia
		colly.Async(true),                        // allow concurrent requests
	)

	// set crawling limits and random delays to avoid hitting Wikipedia too quickly
	c.Limit(&colly.LimitRule{
		Parallelism: 2,               // only two concurrent requests at a time
		RandomDelay: 2 * time.Second, // wait up to 2 seconds randomly between requests
	})

	// before each request print which URL we’re about to visit
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting:", r.URL)
	})

	// when the crawler finds a <div>
	c.OnHTML("div.mw-parser-output", func(e *colly.HTMLElement) {
		textParts := []string{} // holds each paragraph’s text

		// for each <p> tag inside the content div, extract and clean its text
		e.ForEach("h1, h2, h3, h4, h5, h6, p", func(_ int, el *colly.HTMLElement) {
			t := strings.TrimSpace(el.Text) // remove extra whitespace/newlines
			if t != "" {                    // skip empty paragraphs
				textParts = append(textParts, t)

			}
		})

		// skip empty divs
		if len(textParts) == 0 {
			return
		}

		cleanText := strings.Join(textParts, "\n")

		// remove LaTeX-like markup such as \mathcal{...}, \displaystyle, etc.
		latexPattern := regexp.MustCompile(`\\[a-zA-Z]+(\{[^}]*\})?`)
		cleanText = latexPattern.ReplaceAllString(cleanText, "")

		// collapse excessive whitespace and newlines
		spacePattern := regexp.MustCompile(`\s+`)
		cleanText = spacePattern.ReplaceAllString(cleanText, " ")

		// build a Page
		page := Page{
			URL:       e.Request.URL.String(),
			Text:      cleanText, // join paragraphs with newlines
			CrawledAt: time.Now(),
		}

		// convert the struct into JSON format
		b, _ := json.Marshal(page)

		// write the JSON object as one line to the output file
		file.Write(b)
		file.Write([]byte("\n")) // newline between records
	})

	// measure runtime
	start := time.Now()
	for _, u := range urls {
		c.Visit(u)
	}
	c.Wait()

	fmt.Printf("Crawl completed in %v. Output saved to %s\n", time.Since(start), outputFile)
}

// read URLs line by line from the text file
func readURLs(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			urls = append(urls, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return urls, nil
}
