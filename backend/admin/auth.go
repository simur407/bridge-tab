package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	sessionCookie = "admin_session"
	sessionHours  = 12
	loginTitle    = "Logowanie"
	badPassword   = "Nieprawidłowe hasło"
	badCSRF       = "nieprawidłowy token"
)

func passwordOK(given string) bool {
	sumGiven := sha256.Sum256([]byte(given))
	sumExpected := sha256.Sum256([]byte(adminPassword))
	return subtle.ConstantTimeCompare(sumGiven[:], sumExpected[:]) == 1
}

func issueSession() string {
	nonce := make([]byte, 16)
	if _, err := rand.Read(nonce); err != nil {
		panic(err)
	}
	expiry := time.Now().Add(sessionHours * time.Hour).UTC().Unix()
	payload := strconv.FormatInt(expiry, 10) + "." + hex.EncodeToString(nonce)
	return payload + "." + hex.EncodeToString(sign(payload))
}

func verifySession(value string) (string, bool) {
	parts := strings.Split(value, ".")
	if len(parts) != 3 {
		return "", false
	}
	payload := parts[0] + "." + parts[1]
	got, err := hex.DecodeString(parts[2])
	if err != nil {
		return "", false
	}
	if !hmac.Equal(got, sign(payload)) {
		return "", false
	}
	expiry, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || time.Now().Unix() > expiry {
		return "", false
	}
	return parts[1], true
}

func csrfFor(nonce string) string {
	return hex.EncodeToString(sign("csrf." + nonce))
}

func csrfOK(got, nonce string) bool {
	decoded, err := hex.DecodeString(got)
	if err != nil {
		return false
	}
	return hmac.Equal(decoded, sign("csrf."+nonce))
}

func sign(payload string) []byte {
	mac := hmac.New(sha256.New, []byte(adminPassword))
	mac.Write([]byte(payload))
	return mac.Sum(nil)
}

func readSession(c *fiber.Ctx) (string, bool) {
	return verifySession(c.Cookies(sessionCookie))
}

func csrfToken(c *fiber.Ctx) (string, error) {
	nonce, _ := c.Locals("nonce").(string)
	if nonce == "" {
		return "", errors.New("no session")
	}
	return csrfFor(nonce), nil
}

func checkCSRF(c *fiber.Ctx) error {
	nonce, _ := c.Locals("nonce").(string)
	if !csrfOK(c.FormValue("csrf"), nonce) {
		return errors.New(badCSRF)
	}
	return nil
}

func requireAdmin(c *fiber.Ctx) error {
	if c.Path() == "/login" || strings.HasPrefix(c.Path(), "/css/") {
		return c.Next()
	}
	nonce, ok := readSession(c)
	if !ok {
		return c.Redirect("/login", fiber.StatusSeeOther)
	}
	c.Locals("nonce", nonce)
	return c.Next()
}

func setSessionCookie(c *fiber.Ctx, value string) {
	c.Cookie(&fiber.Cookie{
		Name:     sessionCookie,
		Value:    value,
		HTTPOnly: true,
		SameSite: "Lax",
		Secure:   cookieSecure,
		Path:     "/",
		Expires:  time.Now().Add(sessionHours * time.Hour),
	})
}

func clearSessionCookie(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     sessionCookie,
		Value:    "",
		HTTPOnly: true,
		SameSite: "Lax",
		Secure:   cookieSecure,
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})
}

func getLogin(c *fiber.Ctx) error {
	if _, ok := readSession(c); ok {
		return c.Redirect("/tournaments", fiber.StatusSeeOther)
	}
	return view(c, "login", fiber.Map{"Title": loginTitle})
}

func postLogin(c *fiber.Ctx) error {
	if !passwordOK(c.FormValue("password")) {
		return c.Status(fiber.StatusUnauthorized).Render("login", prepare(c, fiber.Map{
			"Title": loginTitle,
			"Error": badPassword,
		}), "layout")
	}
	setSessionCookie(c, issueSession())
	return c.Redirect("/tournaments", fiber.StatusSeeOther)
}

func postLogout(c *fiber.Ctx) error {
	if err := checkCSRF(c); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	clearSessionCookie(c)
	return c.Redirect("/login", fiber.StatusSeeOther)
}
