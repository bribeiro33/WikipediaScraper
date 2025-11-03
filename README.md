# Go Web Crawler/Scraper

This project implements a Go-based web crawler and scraper that collects text data from Wikipedia pages, similar to what a Python/Scrapy crawler would produce.

It uses the Gocolly (https://github.com/gocolly/colly) framework to asynchronously crawl web pages, extract meaningful textual information, and save it in a newline-delimited JSON (NDJSON) format for analysis.

It demonstrates:

* Asynchronous crawling and scraping with Colly

* Clean extraction of semantic text content \(\<h1>, \<h2>, \<h3>, \<h4>, \<p>)

* Data cleaning to remove LaTeX/HTML noise from Wikipedia pages

* JSON Lines output for easy downstream data analysis

* Command-line arguments for flexible input/output handling

* Crawl timing comparison with a Python/Scrapy implementation

Requires Go 1.22 or higher

## Structure
```console
wikipedia_crawler/
├── main.go                        # Main crawler and scraper logic
├── urls.txt                       # Input list of Wikipedia URLs to crawl
└── wikipedia_output.ndjson        # NDJSON file (output of crawl)     

```

## Features
### Web Crawling

* Uses github.com/gocolly/colly/v2 for efficient, concurrent crawling.

* Restricts crawling to the domain en.wikipedia.org.

* Respects polite crawling with:

  * Parallelism limit: 2

  * Random delay between requests

### Text Extraction

* Extracts text from key article sections:

  * Headings

  * Paragraphs

* Cleans output by removing:

  * Escaped LaTeX (\mathcal{}, \displaystyle, etc.)

  * Redundant whitespace and newlines

* Skips irrelevant or empty elements

### Output

* Each page is written as one line of JSON in NDJSON format:
```console
{"url": "https://en.wikipedia.org/wiki/Robot", "text": "Robot\nEtymology\nThe term robot comes from...", "crawled_at": "2025-11-02T17:58:00Z"}
```
* Output file is cutomizable via command-line arg

## How to use
1. Clone repo
2. Initialize go and install dependecies
```console
go mod init wikipedia_crawler
go get github.com/gocolly/colly/v2
go mod tidy
```
4. Create a txt file with your input urls, one URL per line
5. Run program following this format:
```console
go run main.go <input_urls.txt> <output_file_name.ndjson>
```
the output file will automatically be created.
6. Read the output file

## Performance Comparison to Python/Scrapy

All done over the same wifi with an M2 MacBook

Scrapy overall time: 'elapsed_time_seconds': 10.330411

With each website's crawl ranging from a delay of 718 ms to 52 ms.

Go Scraper overall time: ranging from 466 ms to 431 ms.

The time it takes for Scrapy to crawl one page, Go can complete its whole suite. 
