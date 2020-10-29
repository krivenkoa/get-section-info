package errs

import "errors"

var (
	// ErrUnsupportedMedia represents error for unsupported media type
	ErrUnsupportedMedia = errors.New("unsupported http media type")

	// ErrBadRequest represents bad request error
	ErrBadRequest = errors.New("bad request")

	// ErrConvertObjectID represents error for get mongodb object id
	ErrConvertObjectID = errors.New("failed convert object id")

	// ErrNotFound represents entity not found error
	ErrNotFound = errors.New("not found")

	// ErrInternalDatabase represents internal database error
	ErrInternalDatabase = errors.New("internal database error")
)
