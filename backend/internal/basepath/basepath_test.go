package basepath

import "testing"

func TestNormalize(t *testing.T) {
	cases := map[string]string{
		"":           "",
		"/":          "",
		"  ":         "",
		"admin":      "/admin",
		"/admin":     "/admin",
		"/admin/":    "/admin",
		" /tablet/ ": "/tablet",
	}
	for in, want := range cases {
		if got := Normalize(in); got != want {
			t.Fatalf("Normalize(%q)=%q want %q", in, got, want)
		}
	}
}

func TestJoin(t *testing.T) {
	cases := []struct {
		base, path, want string
	}{
		{"", "/login", "/login"},
		{"", "/", "/"},
		{"", "login", "/login"},
		{"/admin", "/login", "/admin/login"},
		{"/admin", "/", "/admin"},
		{"/admin", "tournaments", "/admin/tournaments"},
		{"/admin", "/tournaments/1?x=1", "/admin/tournaments/1?x=1"},
	}
	for _, tc := range cases {
		if got := Join(tc.base, tc.path); got != tc.want {
			t.Fatalf("Join(%q,%q)=%q want %q", tc.base, tc.path, got, tc.want)
		}
	}
}

func TestCookiePath(t *testing.T) {
	if got := CookiePath(""); got != "/" {
		t.Fatalf("empty: %q", got)
	}
	if got := CookiePath("/admin"); got != "/admin" {
		t.Fatalf("admin: %q", got)
	}
}
