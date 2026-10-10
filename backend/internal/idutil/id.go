package idutil

import "strings"

// SameId reports whether two identifiers are the same UUID.
// Postgres returns uuid values in lowercase, while callers may supply uppercase.
func SameId[T ~string](a, b T) bool {
	return strings.EqualFold(string(a), string(b))
}
