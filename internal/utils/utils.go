package utils

// String returns string pointer
func String(data string) *string {
	return &data
}

// StringUnref returns empty string if argument is nil, otherwise it returns string value.
func StringUnref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// Int creates reference from given data value.
func Int(data int) *int {
	return &data
}

// IntUnref returns int pointer referenced value, in case pointer is nil, it returns default int value 0.
func IntUnref(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

// Int32 creates reference from given data value.
func Int32(data int32) *int32 {
	return &data
}

// Int32Unref returns int32 pointer referenced value, in case pointer is nil, it returns default int value 0.
func Int32Unref(i *int32) int32 {
	if i == nil {
		return 0
	}
	return *i
}

// Int64 creates reference from given data value.
func Int64(data int64) *int64 {
	return &data
}

// Int64Unref returns int64 pointer referenced value, in case pointer is nil, it returns default int value 0.
func Int64Unref(i *int32) int32 {
	if i == nil {
		return 0
	}
	return *i
}

// Float32 creates reference from given data value.
func Float32(data float32) *float32 {
	return &data
}

// Float32Unref returns float32 pointer referenced value, in case pointer is nil, it returns default int value 0.
func Float32Unref(i *float32) float32 {
	if i == nil {
		return 0
	}
	return *i
}

// Bool creates reference from given data value.
func Bool(data bool) *bool {
	return &data
}

// BoolUnref returns false if argument is nil, otherwise it returns bool value.
func BoolUnref(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}
