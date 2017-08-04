/**
 * This file is for providing to get web page from weblio.jp
 */

package main

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
)

// getWeblioDict function prints dictionary in weblio.jp.
// Top of discription block and wikipedia-jp discription will be print.
func getWeblioDict(word string) {
	var str string
	yellow := color.New(color.Bold, color.FgYellow).SprintFunc()

	// get page from weblio.jp
	url := "http://www.weblio.jp/content/" + word
	doc, _ := goquery.NewDocument(url)

	// trim elemets
	doc.Find(".NetDicHead .midashigo > span").Empty()

	// print Midashigo
	str = doc.Find(".NetDicHead .midashigo").Text()
	if str != "" {
		fmt.Println(yellow(str))
	} else {
		str = word
		fmt.Println(yellow(str))
	}

	// print dictionary bodies
	doc.Find("#cont .kiji .NetDicBody > div > div > div > div").Each(func(_ int, s *goquery.Selection) {
		str = strings.TrimSpace(s.Text())
		if str != "" {
			fmt.Println(str)
		}
	})
	// print Wikipedia-ja text
	str = strings.TrimSpace(doc.Find(".kiji .Wkpja .WkpjaTs").Next().Text())
	if str != "" {
		fmt.Println(str)
	}
}

// getWeblioEijiDict function prints dictionary in ejje.weblio.jp.
// summary and example writing will be print.
func getWeblioEijiDict(word string) {
	var str string
	green := color.New(color.Bold, color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	bold := color.New(color.Bold).SprintFunc()

	// get page from ejje.weblio.jp
	url := "http://ejje.weblio.jp/content/" + word
	doc, _ := goquery.NewDocument(url)

	// print query
	str = doc.Find("#h1Query").Text()
	if str != "" {
		fmt.Println(green(str))
	} else {
		str = word
		fmt.Println(green(str))
	}

	// print summary
	summary := doc.Find(".summaryM table tbody tr td.content-explanation").Text()
	fmt.Println(summary)

	// print example documents
	fmt.Print(bold("\n[Examples]\n"))
	doc.Find(".qotC").Each(func(_ int, s *goquery.Selection) {
		s.Find(".qotCE > span").Empty()
		s.Find(".qotCJE > span").Empty()
		s.Find(".qotCJJ > span").Empty()
		exampleHead := s.Find(".qotCE").Text()
		exampleTail := s.Find(".qotCJ").Text()
		if exampleHead == "" {
			exampleHead = s.Find(".qotCJE").Text()
			exampleTail = s.Find(".qotCJJ").Text()
		}
		if exampleHead != "" {
			fmt.Println(yellow(exampleHead))
			fmt.Println(exampleTail)
		}
	})

	// print more link
	moreExample, _ := doc.Find(".hlt_SNTCE .kiji > a").Attr("href")
	ul := color.New(color.Underline).SprintFunc()
	fmt.Println("\n", ul(moreExample))
	fmt.Println("")
}
