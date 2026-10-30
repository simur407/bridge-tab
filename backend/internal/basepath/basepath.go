package basepath

import "strings"

// Normalize turns an optional env value into a base path with a leading slash
// and no trailing slash. Empty, "/", and whitespace mean "served at root".
func Normalize(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "/" {
		return ""
	}
	if !strings.HasPrefix(raw, "/") {
		raw = "/" + raw
	}
	return strings.TrimRight(raw, "/")
}

// Join prefixes an absolute app path with base. path should start with "/".
// When base is empty, path is returned unchanged (aside from ensuring a leading "/").
func Join(base, path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		path = "/"
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if strings.HasPrefix(path, "//") {
		path = "/"
	}
	if base == "" {
		return path
	}
	if path == "/" {
		return base
	}
	return base + path
}

// CookiePath is the Path attribute for cookies scoped to this app.
func CookiePath(base string) string {
	if base == "" {
		return "/"
	}
	return base
}
