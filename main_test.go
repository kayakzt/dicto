package main

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
)

func TestRunWithDepsInvalidArgs(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := runWithDeps([]string{}, &stdout, &stderr, failIfCalled, failIfCalled)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "invalid arguments") {
		t.Fatalf("expected invalid arguments error, got %v", err)
	}
	if !strings.Contains(stderr.String(), "Usage: dicto") {
		t.Fatalf("expected usage on stderr, got %q", stderr.String())
	}
}

func TestRunWithDepsEiji(t *testing.T) {
	var stdout, stderr bytes.Buffer
	var gotWord string

	eiji := func(word string, out io.Writer) error {
		gotWord = word
		return nil
	}

	err := runWithDeps([]string{"-e", "hello"}, &stdout, &stderr, failIfCalled, eiji)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotWord != "hello" {
		t.Fatalf("expected word hello, got %q", gotWord)
	}
}

func TestRunWithDepsDefault(t *testing.T) {
	var stdout, stderr bytes.Buffer
	var gotWord string

	dict := func(word string, out io.Writer) error {
		gotWord = word
		return nil
	}

	err := runWithDeps([]string{"world"}, &stdout, &stderr, dict, failIfCalled)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotWord != "world" {
		t.Fatalf("expected word world, got %q", gotWord)
	}
}

func failIfCalled(word string, out io.Writer) error {
	return errors.New("unexpected call")
}
