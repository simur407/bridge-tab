package main

import (
	"errors"
	"strconv"
	"strings"

	tournament_cmd "bridge-tab/internal/tournament-management/application/command"
	tournament_domain "bridge-tab/internal/tournament-management/domain"

	"github.com/gofiber/fiber/v2"
)

func getBoardCards(c *fiber.Ctx) error {
	play, tournament, a, err := loadTablePlayForDeal(c)
	if err != nil {
		if redirect, ok := err.(redirectError); ok {
			return c.Redirect(path(redirect.To), fiber.StatusSeeOther)
		}
		return err
	}
	dealNo, err := strconv.Atoi(c.Params("dealNo"))
	if err != nil || dealNo < 1 {
		return fiber.NewError(fiber.StatusNotFound, "nie znaleziono rozdania")
	}
	board, ok := findBoard(play.Boards, dealNo)
	if !ok {
		return fiber.NewError(fiber.StatusNotFound, "to rozdanie nie jest grane w tej rundzie")
	}
	if !tournament.ScoreEntryMode.TabletCanRegister() {
		return c.Redirect(path("/play"), fiber.StatusSeeOther)
	}
	complete, err := boardHandsComplete(c, a.TournamentId, dealNo)
	if err != nil {
		return err
	}
	if complete {
		return c.Redirect(path("/play/boards/"+strconv.Itoa(dealNo)), fiber.StatusSeeOther)
	}
	data := boardPageData(tournament, a, play, board)
	data["CardsPostURL"] = path("/play/boards/" + strconv.Itoa(dealNo) + "/cards")
	return playView(c, "board-cards", data)
}

func postBoardCards(c *fiber.Ctx) error {
	play, tournament, a, err := loadTablePlayForDeal(c)
	if err != nil {
		if redirect, ok := err.(redirectError); ok {
			return c.Redirect(path(redirect.To), fiber.StatusSeeOther)
		}
		return err
	}
	dealNo, err := strconv.Atoi(c.Params("dealNo"))
	if err != nil || dealNo < 1 {
		return fiber.NewError(fiber.StatusNotFound, "nie znaleziono rozdania")
	}
	board, ok := findBoard(play.Boards, dealNo)
	if !ok {
		return fiber.NewError(fiber.StatusNotFound, "to rozdanie nie jest grane w tej rundzie")
	}
	if !tournament.ScoreEntryMode.TabletCanRegister() {
		return c.Redirect(path("/play"), fiber.StatusSeeOther)
	}

	s := newStores(c)
	cmd := tournament_cmd.SetBoardHandsCommand{
		TournamentId: a.TournamentId,
		BoardNo:      dealNo,
		North:        strings.TrimSpace(c.FormValue("north")),
		East:         strings.TrimSpace(c.FormValue("east")),
		South:        strings.TrimSpace(c.FormValue("south")),
		West:         strings.TrimSpace(c.FormValue("west")),
		Overwrite:    false,
	}
	err = cmd.Execute(s.tournament)
	if err != nil && !errors.Is(err, tournament_domain.ErrBoardHandsAlreadySet) {
		data := boardPageData(tournament, a, play, board)
		data["CardsPostURL"] = path("/play/boards/" + strconv.Itoa(dealNo) + "/cards")
		data["Error"] = explain(err).Error()
		return c.Status(fiber.StatusBadRequest).Render("board-cards", prepare(c, data), "play-layout")
	}
	return c.Redirect(path("/play/boards/"+strconv.Itoa(dealNo)), fiber.StatusSeeOther)
}

func boardHandsComplete(c *fiber.Ctx, tournamentId string, dealNo int) (bool, error) {
	s := newStores(c)
	protocol, err := s.boards.FindByBoardNo(&tournamentId, dealNo)
	if err != nil {
		return false, err
	}
	if protocol == nil {
		return false, nil
	}
	return protocol.HandsComplete(), nil
}

func redirectAfterTabletAdd(c *fiber.Ctx, tournamentId string, dealNo int) error {
	complete, err := boardHandsComplete(c, tournamentId, dealNo)
	if err != nil {
		complete = true
	}
	target := "/play/boards/" + strconv.Itoa(dealNo)
	if !complete {
		target = target + "/cards"
	}
	c.Set("HX-Redirect", path(target))
	return c.SendStatus(fiber.StatusNoContent)
}
