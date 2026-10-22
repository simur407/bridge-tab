package main

import (
	"database/sql"
	"time"

	umpirecall "bridge-tab/internal/umpire-call"

	"github.com/gofiber/fiber/v2"
)

func getUmpireCall(c *fiber.Ctx) error {
	play, _, a, err := loadTablePlay(c)
	if err != nil {
		if redirect, ok := err.(redirectError); ok {
			if isHTMX(c) {
				c.Set("HX-Redirect", redirect.To)
				return c.SendStatus(fiber.StatusNoContent)
			}
			return c.Redirect(redirect.To, fiber.StatusSeeOther)
		}
		return err
	}
	if play.Empty || play.RoundNo == 0 {
		return c.SendString("")
	}
	return renderUmpireCall(c, a.TournamentId, a.TableNumber)
}

func postUmpireCall(c *fiber.Ctx) error {
	play, _, a, err := loadTablePlay(c)
	if err != nil {
		if redirect, ok := err.(redirectError); ok {
			if isHTMX(c) {
				c.Set("HX-Redirect", redirect.To)
				return c.SendStatus(fiber.StatusNoContent)
			}
			return c.Redirect(redirect.To, fiber.StatusSeeOther)
		}
		return err
	}
	if play.Empty || play.RoundNo == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "ten stół nie ma rundy")
	}
	_, err = umpireRepo(c).Insert(umpirecall.Call{
		TournamentId: a.TournamentId,
		TableNumber:  a.TableNumber,
		RoundNo:      play.RoundNo,
		NsTeamNumber: play.NsTeamNumber,
		EwTeamNumber: play.EwTeamNumber,
		Source:       umpirecall.SourceTablet,
	})
	if err != nil {
		return err
	}
	if isHTMX(c) {
		return renderUmpireCall(c, a.TournamentId, a.TableNumber)
	}
	return c.Redirect("/play", fiber.StatusSeeOther)
}

func renderUmpireCall(c *fiber.Ctx, tournamentId string, tableNumber int) error {
	call, err := umpireRepo(c).LatestForTable(tournamentId, tableNumber)
	if err != nil {
		return err
	}
	return fragment(c, "umpire-call", umpireButtonData(call, "/play/umpire-call"))
}

func umpireRepo(c *fiber.Ctx) *umpirecall.Repository {
	return &umpirecall.Repository{
		Ctx: c.UserContext(),
		Tx:  c.Locals("tx").(*sql.Tx),
	}
}

func umpireButtonData(call *umpirecall.Call, url string) fiber.Map {
	state, left := umpirecall.ButtonFrom(call, time.Now())
	label := "Wezwij sędziego"
	disabled := false
	coming := false
	var enableInMs int64
	var seconds int
	switch state {
	case umpirecall.ButtonWaiting:
		label = "Wezwanie wysłane"
		disabled = true
	case umpirecall.ButtonComing:
		label = "Ponowne wezwanie"
		disabled = true
		coming = true
		enableInMs = left.Milliseconds()
		if enableInMs < 1 {
			enableInMs = 1
		}
		seconds = int((left + time.Second - 1) / time.Second)
		if seconds < 1 {
			seconds = 1
		}
	}
	return fiber.Map{
		"Label":      label,
		"Disabled":   disabled,
		"Coming":     coming,
		"Waiting":    state == umpirecall.ButtonWaiting,
		"Polling":    state == umpirecall.ButtonWaiting,
		"State":      string(state),
		"EnableInMs": enableInMs,
		"Seconds":    seconds,
		"PollURL":    url,
		"PostURL":    url,
	}
}
