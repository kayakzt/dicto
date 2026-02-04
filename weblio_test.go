package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
)

func TestRenderWeblioDict(t *testing.T) {
	doc := mustDocument(t, weblioDictHTML)
	var out bytes.Buffer

	if err := renderWeblioDict(doc, "fallback", &out, fmt.Sprint); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	gotLines := strings.Split(strings.TrimSpace(out.String()), "\n")
	wantLines := []string{"Head", "Meaning 1", "Wiki text"}

	if len(gotLines) != len(wantLines) {
		t.Fatalf("expected %d lines, got %d: %q", len(wantLines), len(gotLines), gotLines)
	}
	for i, want := range wantLines {
		if gotLines[i] != want {
			t.Fatalf("line %d: expected %q, got %q", i+1, want, gotLines[i])
		}
	}
}

func TestRenderWeblioDictFallbackHeading(t *testing.T) {
	doc := mustDocument(t, weblioDictNoHeadHTML)
	var out bytes.Buffer

	if err := renderWeblioDict(doc, "fallback", &out, fmt.Sprint); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	gotLines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if gotLines[0] != "fallback" {
		t.Fatalf("expected fallback heading, got %q", gotLines[0])
	}
}

func TestRenderWeblioEijiDict(t *testing.T) {
	origNoColor := color.NoColor
	color.NoColor = true
	t.Cleanup(func() { color.NoColor = origNoColor })

	doc := mustDocument(t, weblioEijiHTML)
	var out bytes.Buffer

	if err := renderWeblioEijiDict(doc, "fallback", &out, fmt.Sprint, fmt.Sprint, fmt.Sprint); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := out.String()
	assertInOrder(t, got, []string{
		"Query",
		"Summary",
		"[Examples]",
		"Example head 1",
		"Example tail 1",
		"Example head 2",
		"Example tail 2",
		"https://example.com/more",
	})
}

func mustDocument(t *testing.T, html string) *goquery.Document {
	t.Helper()

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatalf("failed to parse html: %v", err)
	}
	return doc
}

func assertInOrder(t *testing.T, s string, parts []string) {
	t.Helper()

	idx := 0
	for _, part := range parts {
		pos := strings.Index(s[idx:], part)
		if pos < 0 {
			t.Fatalf("missing %q in output: %q", part, s)
		}
		idx += pos + len(part)
	}
}

const weblioDictHTML = `
<html><body>
  <div class="NetDicHead"><div class="midashigo"><span>ignore</span> Head </div></div>
  <div id="cont">
    <div class="kiji">
      <div class="NetDicBody"><div><div><div><div> Meaning 1 </div></div></div></div></div>
    </div>
  </div>
  <div class="kiji"><div class="Wkpja"><div class="WkpjaTs">Title</div><div> Wiki text </div></div></div>
</body></html>
`

const weblioDictNoHeadHTML = `
<html><body>
  <div class="NetDicHead"><div class="midashigo"><span>ignore</span> </div></div>
  <div id="cont">
    <div class="kiji">
      <div class="NetDicBody"><div><div><div><div> Meaning 1 </div></div></div></div></div>
    </div>
  </div>
</body></html>
`

const weblioEijiHTML = `
<html><body>
  <div id="h1Query">Query</div>
  <div class="summaryM">
    <p><span class="content-explanation"> Summary </span></p>
  </div>
  <div class="qotC">
    <div class="qotCE"><span>rm</span> Example head 1 </div>
    <div class="qotCJ"> Example tail 1 </div>
  </div>
  <div class="qotC">
    <div class="qotCJE"><span>rm</span> Example head 2 </div>
    <div class="qotCJJ"> Example tail 2 </div>
  </div>
  <div class="hlt_SNTCE"><div class="kiji"><a href="https://example.com/more">more</a></div></div>
</body></html>
`
