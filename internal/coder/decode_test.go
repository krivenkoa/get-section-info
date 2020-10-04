package coder

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestReadBodyDefaultJSON(t *testing.T) {
	req := TestRequest{}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/tests", strings.NewReader(`{"id":"onotole","value":"upyachka"}`))
	err := ReadBody(w, r, &req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if req.ID != "onotole" || req.Value != "upyachka" {
		t.Fatalf("wrong result: %v", req)
	}
}

func TestReadBodyJSON(t *testing.T) {
	req := TestRequest{}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/tests", strings.NewReader(`{"id":"onotole","value":"upyachka"}`))
	r.Header.Set("Content-Type", "application/json")
	err := ReadBody(w, r, &req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if req.ID != "onotole" || req.Value != "upyachka" {
		t.Fatalf("wrong result: %v", req)
	}
}

func TestReadBodyComplicatedJSONMime(t *testing.T) {
	req := TestRequest{}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/tests", strings.NewReader(`{"id":"onotole","value":"upyachka"}`))
	r.Header.Set("Content-Type", "application/json;charset=utf-8")
	err := ReadBody(w, r, &req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if req.ID != "onotole" || req.Value != "upyachka" {
		t.Fatalf("wrong result: %v", req)
	}
}

func TestReadBodyBadContentType(t *testing.T) {
	req := TestRequest{}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/tests", strings.NewReader(`{"id":"onotole","value":"upyachka"}`))
	r.Header.Set("Content-Type", "application!json")
	err := ReadBody(w, r, &req)
	if err == nil {
		t.Fatalf("expected error")
	}

	if w.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("bad response code %v", w.Code)
	}
}
