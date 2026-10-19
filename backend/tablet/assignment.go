package main

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	tournamentCookie = "tablet_tournament"
	tableCookie      = "tablet_table"
	holdRoundCookie  = "tablet_hold_round"
	assignmentDays   = 30
)

type assignment struct {
	TournamentId string
	TableNumber  int
}

func readAssignment(c *fiber.Ctx) (assignment, bool) {
	tournamentId := c.Cookies(tournamentCookie)
	tableRaw := c.Cookies(tableCookie)
	if tournamentId == "" || tableRaw == "" {
		return assignment{}, false
	}
	tableNumber, err := strconv.Atoi(tableRaw)
	if err != nil || tableNumber < 1 {
		return assignment{}, false
	}
	return assignment{TournamentId: tournamentId, TableNumber: tableNumber}, true
}

func setAssignment(c *fiber.Ctx, a assignment) {
	expires := time.Now().Add(assignmentDays * 24 * time.Hour)
	c.Cookie(&fiber.Cookie{
		Name:     tournamentCookie,
		Value:    a.TournamentId,
		HTTPOnly: true,
		SameSite: "Lax",
		Secure:   cookieSecure,
		Path:     "/",
		Expires:  expires,
	})
	c.Cookie(&fiber.Cookie{
		Name:     tableCookie,
		Value:    strconv.Itoa(a.TableNumber),
		HTTPOnly: true,
		SameSite: "Lax",
		Secure:   cookieSecure,
		Path:     "/",
		Expires:  expires,
	})
}

func readHoldRound(c *fiber.Ctx) int {
	raw := c.Cookies(holdRoundCookie)
	if raw == "" {
		return 0
	}
	round, err := strconv.Atoi(raw)
	if err != nil || round < 1 {
		return 0
	}
	return round
}

func setHoldRound(c *fiber.Ctx, roundNo int) {
	c.Cookie(&fiber.Cookie{
		Name:     holdRoundCookie,
		Value:    strconv.Itoa(roundNo),
		HTTPOnly: true,
		SameSite: "Lax",
		Secure:   cookieSecure,
		Path:     "/",
		Expires:  time.Now().Add(assignmentDays * 24 * time.Hour),
	})
}

func clearHoldRound(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     holdRoundCookie,
		Value:    "",
		HTTPOnly: true,
		SameSite: "Lax",
		Secure:   cookieSecure,
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})
}

func assignmentLocked(a assignment, startedAt string) bool {
	return a.TournamentId != "" && startedAt != ""
}

func isHTMX(c *fiber.Ctx) bool {
	return c.Get("HX-Request") == "true"
}
