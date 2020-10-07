package coder

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	chimid "github.com/go-chi/chi/middleware"

	"skat-vending.com/selection-info/internal/errs"
	"skat-vending.com/selection-info/pkg/api"
)

func TestWriteJSON(t *testing.T) {
	res := TestResponse{
		ID:    "onotole",
		Value: "upyachka",
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/tests", nil)
	r.Header.Set("Accept", "application/json")
	err := WriteData(w, r, &res, http.StatusOK)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	bytes, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if contentType := w.Header().Get("Content-Type"); contentType != "application/json" {
		t.Fatalf("wrong content-type %v", w.Header().Get("Content-Type"))
	}

	if w.Header().Get("X-Request-ID") != chimid.GetReqID(r.Context()) {
		t.Fatalf("req_id mismatch")
	}

	control := TestResponse{}
	err = json.Unmarshal(bytes, &control)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(control, res) {
		t.Fatalf("expected: %v got: %v", res, control)
	}
}

func TestWriteError(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/tests", nil)
	r.Header.Set("Accept", "application/json")
	err := WriteError(w, r, http.StatusNotFound, errs.ErrUnsupportedMedia.Error())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if contentType := w.Header().Get("Content-Type"); contentType != "application/json" {
		t.Fatalf("wrong content-type %v", w.Header().Get("Content-Type"))
	}

	control := api.HTTPError{}
	err = json.NewDecoder(w.Body).Decode(&control)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if control.Error != errs.ErrUnsupportedMedia.Error() && !control.Status {
		t.Fatalf("got: %v", control)
	}
}
