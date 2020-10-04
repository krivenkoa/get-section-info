package coder

import (
	"net/http"
	"sort"
	"strings"

	"github.com/golang/gddo/httputil/header"
)

const (
	// ContentTypeJSON represents "application/json" content type
	ContentTypeJSON = "application/json"
)

// MatchMimetype returns true if given mimetype matches given mime pattern.
// e.g. application/json matches application/*, but application/json and plain/text doesn't.
func MatchMimetype(mime, pattern string) bool {
	mimeSlice := strings.Split(mime, "/")
	patternSlice := strings.Split(pattern, "/")
	if len(mimeSlice) != len(patternSlice) {
		return false
	}
	for i := 0; i != len(mimeSlice); i++ {
		if patternSlice[i] != "*" && patternSlice[i] != mimeSlice[i] {
			return false
		}
	}
	return true
}

// FindBestMatchMimeType returns best matching supported mimetype for a given header key with all respect to weights.
// Algorithm complexity is O(M*N) where N is a number of supported mimetypes and M is a number of accept header values.
func FindBestMatchMimeType(supportedMimeTypes []string, h http.Header, key string) string {
	accept := header.ParseAccept(h, key)
	sort.SliceStable(accept[:], func(i, j int) bool {
		return accept[i].Q > accept[j].Q
	})

	for _, pattern := range accept {
		for _, mime := range supportedMimeTypes {
			if MatchMimetype(mime, pattern.Value) {
				return mime
			}
		}
	}
	return ""
}
