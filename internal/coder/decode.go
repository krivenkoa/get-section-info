package coder

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"

	wrapper "github.com/pkg/errors"

	"skat-vending.com/selection-info/internal/utils"
	"skat-vending.com/selection-info/pkg/api"
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
			res := &api.Section{
				Success: false,
				Description: []api.Description{
					{
						Message:    "unsupported media type",
						Stacktrace: utils.String(err.Error()),
					},
				},
			}
			WriteBadCode(w, r, res, http.StatusUnsupportedMediaType)
			return fmt.Errorf("couldn't parse content-type: '%s'", contentType)
		}
		contentType = mediatype
	}
	switch contentType {
	case ContentTypeJSON:
		{
			err := json.NewDecoder(r.Body).Decode(req)
			if err != nil {
				res := &api.Section{
					Success: false,
					Description: []api.Description{
						{
							Message:    "error reading content of json type",
							Stacktrace: utils.String(err.Error()),
						},
					},
				}
				WriteBadCode(w, r, res, http.StatusBadRequest)
				return wrapper.Wrap(err, "error reading content of json type")
			}
		}
	default:
		{
			res := &api.Section{
				Success: false,
				Description: []api.Description{
					{
						Message: "unsupported content-type",
						Reason:  utils.String(fmt.Sprintf("content-type: '%s'", contentType)),
					},
				},
			}
			WriteBadCode(w, r, res, http.StatusUnsupportedMediaType)
			return fmt.Errorf("unsupported content-type: '%s'", contentType)
		}
	}
	return nil
}
