package main

import "testing"

func TestPasswordMatch(t *testing.T) {
	adminPassword = "secret"
	if !passwordOK("secret") {
		t.Fatal("expected matching password")
	}
	if passwordOK("secret ") || passwordOK("secre") || passwordOK("") {
		t.Fatal("expected password mismatch")
	}
}

func TestSessionRoundTrip(t *testing.T) {
	adminPassword = "test-secret"
	value := issueSession()
	nonce, ok := verifySession(value)
	if !ok || nonce == "" {
		t.Fatal("session was rejected")
	}
	if _, ok := verifySession(value + "aa"); ok {
		t.Fatal("tampered session was accepted")
	}
	if !csrfOK(csrfFor(nonce), nonce) {
		t.Fatal("csrf token does not match the session")
	}
	if csrfOK(csrfFor(nonce), "other") {
		t.Fatal("csrf token matched a different session")
	}
}
