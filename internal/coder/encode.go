package coder

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/middleware"
	wrapper "github.com/pkg/errors"

	"skat-vending.com/selection-info/internal/errs"
)

// defaultAccept
const defaultAccept = ContentTypeJSON

// WriteError writes error to response
func WriteError(w http.ResponseWriter, r *http.Request, body interface{}, statusCode int) error {
	inferred := inferAccept(r)
	return wrapper.Wrap(writeData(w, body, statusCode, inferred, middleware.GetReqID(r.Context())), "writing errors")
}

// WriteBadCode writes bad request response code
func WriteBadCode(w http.ResponseWriter, r *http.Request, body interface{}, statusCode int) error {
	return wrapper.Wrap(WriteError(w, r, body, statusCode), "error writing bad status code")
}

// WriteData writes response data
func WriteData(w http.ResponseWriter, r *http.Request, body interface{}, responseCode int) error {
	inferred := inferAccept(r)
	if inferred == "" {
		w.Header().Add("Content-Type", defaultAccept)
		WriteError(w, r, body, http.StatusUnsupportedMediaType)
		return errs.ErrUnsupportedMedia
	}

	return wrapper.Wrap(writeData(w, body, responseCode, inferred, middleware.GetReqID(r.Context())), "unsupported mimetype")
}

func inferAccept(r *http.Request) string {
	if r.Header.Get("Accept") == "" {
		return ""
	}

	return FindBestMatchMimeType([]string{ContentTypeJSON}, r.Header, "Accept")
}

func writeData(w http.ResponseWriter, body interface{}, statusCode int, accept, requestID string) error {
	// for responses with no body don't specify content-type
	if body != nil {
		w.Header().Set("Content-Type", accept)
	}
	// but always set X-Request-ID header
	w.Header().Set("X-Request-ID", requestID)
	w.WriteHeader(statusCode)
	if body == nil {
		return nil
	}
	switch accept {
	case ContentTypeJSON:
		{
			return wrapper.Wrap(json.NewEncoder(w).Encode(body), "error encoding to json")
		}
	}
	return nil
}
