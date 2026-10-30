package main

import "bridge-tab/internal/basepath"

var basePath string

func path(p string) string {
	return basepath.Join(basePath, p)
}

func cookiePath() string {
	return basepath.CookiePath(basePath)
}
