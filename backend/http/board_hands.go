package main

import (
	"errors"
	"strconv"
	"strings"

	"bridge-tab/http/middleware"
	tournament_cmd "bridge-tab/internal/tournament-management/application/command"
	tournament_domain "bridge-tab/internal/tournament-management/domain"
	tournament_management_infra "bridge-tab/internal/tournament-management/infrastructure"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func GetTeamBoardCards(c *fiber.Ctx) error {
	id, err := gameSessionID(c)
	if err != nil {
		return notFound(c, "Niepoprawny identyfikator rozgrywki")
	}
	dealNo, err := dealNumber(c)
	if err != nil {
		return notFound(c, "Nie znaleziono rozdania")
	}
	loaded, err := loadTeamMatches(c, id)
	if err != nil {
		log.Debug(err)
		return notFound(c, translateError(err))
	}
	if !loaded.tournament.ScoreEntryMode.ContestantsCanRegister() {
		return notFound(c, scoreEntryDenied(loaded.tournament.ScoreEntryMode))
	}
	play, ok := teamPlayForDeal(loaded.teamNumber, loaded.tournament.FinishedAt != "", loaded.matches, dealNo)
	if !ok {
		return notFound(c, "To rozdanie nie jest grane przez tę drużynę")
	}
	play = labelTeamPlay(play, loaded.teams)
	board, ok := findBoard(play.Boards, dealNo)
	if !ok {
		return notFound(c, "To rozdanie nie jest grane przez tę drużynę")
	}
	complete, err := boardHandsComplete(c, id, dealNo)
	if err != nil {
		log.Debug(err)
		return notFound(c, translateError(err))
	}
	if complete {
		return c.Redirect(path("/round-registration/"+id+"/boards/"+strconv.Itoa(dealNo)), fiber.StatusSeeOther)
	}
	data := boardPageData(id, loaded.tournament, play, board, c.Locals("user").(middleware.UserMetadata).Id)
	data["CardsPostURL"] = path("/round-registration/" + id + "/boards/" + strconv.Itoa(dealNo) + "/cards")
	data["CardsDoneURL"] = path("/round-registration/" + id + "/boards/" + strconv.Itoa(dealNo))
	return c.Render("team-board-cards", data, "layout")
}

func PostTeamBoardCards(c *fiber.Ctx) error {
	id, err := gameSessionID(c)
	if err != nil {
		return notFound(c, "Niepoprawny identyfikator rozgrywki")
	}
	dealNo, err := dealNumber(c)
	if err != nil {
		return notFound(c, "Nie znaleziono rozdania")
	}
	loaded, err := loadTeamMatches(c, id)
	if err != nil {
		log.Debug(err)
		return notFound(c, translateError(err))
	}
	if !loaded.tournament.ScoreEntryMode.ContestantsCanRegister() {
		return notFound(c, scoreEntryDenied(loaded.tournament.ScoreEntryMode))
	}
	if _, ok := teamPlayForDeal(loaded.teamNumber, loaded.tournament.FinishedAt != "", loaded.matches, dealNo); !ok {
		return notFound(c, "To rozdanie nie jest grane przez tę drużynę")
	}

	north := strings.TrimSpace(c.FormValue("north"))
	east := strings.TrimSpace(c.FormValue("east"))
	south := strings.TrimSpace(c.FormValue("south"))
	west := strings.TrimSpace(c.FormValue("west"))
	cmd := tournament_cmd.SetBoardHandsCommand{
		TournamentId: id,
		BoardNo:      dealNo,
		North:        north,
		East:         east,
		South:        south,
		West:         west,
		Overwrite:    false,
	}
	err = cmd.Execute(&tournament_management_infra.PostgresTournamentRepository{
		Ctx: c.UserContext(),
		Tx:  middleware.GetTransaction(c),
	})
	if err != nil && !errors.Is(err, tournament_domain.ErrBoardHandsAlreadySet) {
		log.Debug(err)
		play, _ := teamPlayForDeal(loaded.teamNumber, loaded.tournament.FinishedAt != "", loaded.matches, dealNo)
		play = labelTeamPlay(play, loaded.teams)
		board, _ := findBoard(play.Boards, dealNo)
		data := boardPageData(id, loaded.tournament, play, board, c.Locals("user").(middleware.UserMetadata).Id)
		data["CardsPostURL"] = path("/round-registration/" + id + "/boards/" + strconv.Itoa(dealNo) + "/cards")
		data["CardsDoneURL"] = path("/round-registration/" + id + "/boards/" + strconv.Itoa(dealNo))
		data["Error"] = translateError(err)
		return c.Status(fiber.StatusBadRequest).Render("team-board-cards", data, "layout")
	}
	return c.Redirect(path("/round-registration/"+id+"/boards/"+strconv.Itoa(dealNo)), fiber.StatusSeeOther)
}

func boardHandsComplete(c *fiber.Ctx, tournamentId string, dealNo int) (bool, error) {
	repo := &tournament_management_infra.PostgresBoardProtocolReadRepository{
		Ctx: c.UserContext(),
		Tx:  middleware.GetTransaction(c),
	}
	protocol, err := repo.FindByBoardNo(&tournamentId, dealNo)
	if err != nil {
		return false, err
	}
	if protocol == nil {
		return false, nil
	}
	return protocol.HandsComplete(), nil
}

func redirectAfterBoardAdd(c *fiber.Ctx, tournamentId string, dealNo int, boardPath string) error {
	complete, err := boardHandsComplete(c, tournamentId, dealNo)
	if err != nil {
		log.Debug(err)
		complete = true
	}
	target := boardPath
	if !complete {
		target = boardPath + "/cards"
	}
	c.Set("HX-Redirect", path(target))
	if !isHTMX(c) {
		return c.Redirect(path(target), fiber.StatusSeeOther)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
