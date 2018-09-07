package main

import (
	"testing"
)

func TestParseRoot(t *testing.T) {
	p := "/"
	t.Logf("Parsing %v", p)

	path := NewPath(p)

	if len(path.Path) > 0 || len(path.ID) > 0 {
		t.Fatalf("Path and ID should be both empty but actually Path: %v, ID: %v", path.Path, path.ID)
	}

	if path.HasID() {
		t.Fatalf("This path should have no ID")
	}
}

func TestParseWithID(t *testing.T) {
	p := "/people/1/"
	t.Logf("Parsing %v", p)

	path := NewPath(p)

	if path.Path != "people" {
		t.Fatalf("Path should be 'people' but actually is: %v", path.Path)
	}

	if path.ID != "1" {
		t.Fatalf("ID should be '1' but actually is: %v", path.ID)
	}

	if !path.HasID() {
		t.Fatalf("This path should have ID")
	}
}

func TestParseWithNoID(t *testing.T) {
	p := "/people/"
	t.Logf("Parsing %v", p)

	path := NewPath(p)

	if path.Path != "people" {
		t.Fatalf("Path should be 'people' but actually is: %v", path.Path)
	}

	if path.ID != "" {
		t.Fatalf("ID should be empty but actually is: %v", path.ID)
	}

	if path.HasID() {
		t.Fatalf("This path should have no ID")
	}
}
