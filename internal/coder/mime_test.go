package coder

import (
	"net/http"
	"testing"
)

func TestMatchMimetype(t *testing.T) {
	if !MatchMimetype("text/html", "text/html") {
		t.Fatalf("mismatch")
	}
	if !MatchMimetype("application/json", "application/json") {
		t.Fatalf("mismatch")
	}
	if !MatchMimetype("application/json", "application/*") {
		t.Fatalf("mismatch")
	}
	if !MatchMimetype("application/json", "*/*") {
		t.Fatalf("mismatch")
	}
}

func TestFindBestMatchMimeType(t *testing.T) {
	h1 := http.Header{"Accept": []string{"text/plain;q=0.6", "text/html", "*/*;q=0.8", "application/xhtml+xml", "application/xml;q=0.9"}}
	h2 := http.Header{}
	if FindBestMatchMimeType([]string{}, h1, "Accept") != "" {
		t.Fatalf("mismatch")
	}
	if FindBestMatchMimeType([]string{"text/plain"}, h2, "Accept") != "" {
		t.Fatalf("mismatch")
	}
	if FindBestMatchMimeType([]string{}, h1, "Accept") != "" {
		t.Fatalf("mismatch")
	}
	if FindBestMatchMimeType(nil, h1, "Accept") != "" {
		t.Fatalf("mismatch")
	}
	if FindBestMatchMimeType([]string{"text/plain", "application/json", "text/html"}, h1, "Accept") != "text/html" {
		t.Fatalf("mismatch")
	}
}
