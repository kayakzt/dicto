// This file provides helpers to fetch pages from weblio.jp.

package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
)

const (
	weblioBaseURL     = "https://www.weblio.jp/content/"
	weblioEijiBaseURL = "https://ejje.weblio.jp/content/"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

// getWeblioDict prints a dictionary entry from weblio.jp.
// It prints the heading and the top dictionary/wikipedia excerpts.
func getWeblioDict(word string, out io.Writer) error {
	yellow := color.New(color.Bold, color.FgYellow).SprintFunc()

	// get page from weblio.jp
	doc, err := fetchDocument(context.Background(), weblioBaseURL+word)
	if err != nil {
		return err
	}
	return renderWeblioDict(doc, word, out, yellow)
}

func renderWeblioDict(doc *goquery.Document, word string, out io.Writer, yellow func(a ...any) string) error {
	// trim elements
	doc.Find(".NetDicHead .midashigo > span").Empty()

	// print midashigo
	head := strings.TrimSpace(doc.Find(".NetDicHead .midashigo").Text())
	if head == "" {
		head = word
	}
	fmt.Fprintln(out, yellow(head))

	// print dictionary bodies
	doc.Find("#cont .kiji > div").Each(func(_ int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" {
			fmt.Fprintln(out, text)
		}
	})
	// print Wikipedia-ja text
	wikipedia := strings.TrimSpace(doc.Find(".kiji .Wkpja .WkpjaTs").Next().Text())
	if wikipedia != "" {
		fmt.Fprintln(out, wikipedia)
	}
	return nil
}

// getWeblioEijiDict prints a dictionary entry from ejje.weblio.jp.
// It prints the summary and example sentences.
func getWeblioEijiDict(word string, out io.Writer) error {
	green := color.New(color.Bold, color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	bold := color.New(color.Bold).SprintFunc()

	// get page from ejje.weblio.jp
	doc, err := fetchDocument(context.Background(), weblioEijiBaseURL+word)
	if err != nil {
		return err
	}
	return renderWeblioEijiDict(doc, word, out, green, yellow, bold)
}

func renderWeblioEijiDict(
	doc *goquery.Document,
	word string,
	out io.Writer,
	green func(a ...any) string,
	yellow func(a ...any) string,
	bold func(a ...any) string,
) error {
	// print query
	query := strings.TrimSpace(doc.Find("#h1Query").Text())
	if query == "" {
		query = word
	}
	fmt.Fprintln(out, green(query))

	// print summary
	summary := strings.TrimSpace(doc.Find(".summaryM table tbody tr td.content-explanation").Text())
	fmt.Fprintln(out, summary)

	// print example documents
	fmt.Fprint(out, bold("\n[Examples]\n"))
	doc.Find(".qotC").Each(func(_ int, s *goquery.Selection) {
		s.Find(".qotCE > span").Empty()
		s.Find(".qotCJE > span").Empty()
		s.Find(".qotCJJ > span").Empty()
		exampleHead := strings.TrimSpace(s.Find(".qotCE").Text())
		exampleTail := strings.TrimSpace(s.Find(".qotCJ").Text())
		if exampleHead == "" {
			exampleHead = strings.TrimSpace(s.Find(".qotCJE").Text())
			exampleTail = strings.TrimSpace(s.Find(".qotCJJ").Text())
		}
		if exampleHead != "" {
			fmt.Fprintln(out, yellow(exampleHead))
			fmt.Fprintln(out, exampleTail)
		}
	})

	// print more link
	moreExample, _ := doc.Find(".hlt_SNTCE .kiji > a").Attr("href")
	ul := color.New(color.Underline).SprintFunc()
	fmt.Fprintln(out, "\n", ul(moreExample))
	fmt.Fprintln(out, "")
	return nil
}

func fetchDocument(ctx context.Context, url string) (*goquery.Document, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Printf("failed to close response body: %v", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	return goquery.NewDocumentFromReader(resp.Body)
}
