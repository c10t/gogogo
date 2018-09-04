package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var mockHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello HTTP Test")
})

func TestOpenVars(t *testing.T) {
	t.Logf("vars state is: %v", vars)

	mockServer := httptest.NewServer(mockHandler)
	defer mockServer.Close()

	req, err := http.NewRequest("GET", mockServer.URL, nil)
	if err != nil {
		t.Fatalf("Error when http.NewRequest: %v", err)
	}

	OpenVars(req)
	t.Logf("vars state is: %v", vars)
	if vars == nil {
		t.Fatalf("vars must not be nil after OpenVars()")
	}
}
