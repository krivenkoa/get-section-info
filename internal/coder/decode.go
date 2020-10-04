package coder

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"

	wrapper "github.com/pkg/errors"
)

const defaultContentType = ContentTypeJSON

// ReadBody reads request body
func ReadBody(w http.ResponseWriter, r *http.Request, req interface{}) error {
	var contentType string
	h := r.Header.Get("Content-Type")
	if h == "" {
		contentType = defaultContentType
	} else {
		// TODO: use second argument (mime params)
		mediatype, _, err := mime.ParseMediaType(h)
		if err != nil {
			WriteBadCode(w, r, http.StatusUnsupportedMediaType, err.Error())
			return fmt.Errorf("couldn't parse content-type: '%s'", contentType)
		}
		contentType = mediatype
	}
	switch contentType {
	case ContentTypeJSON:
		{
			err := json.NewDecoder(r.Body).Decode(req)
			if err != nil {
				WriteBadCode(w, r, http.StatusBadRequest, err.Error())
				return wrapper.Wrap(err, "error reading content of json type")
			}
		}
	default:
		{
			WriteBadCode(w, r, http.StatusUnsupportedMediaType, "unsupported content-type")
			return fmt.Errorf("unsupported content-type: '%s'", contentType)
		}
	}
	return nil
}
