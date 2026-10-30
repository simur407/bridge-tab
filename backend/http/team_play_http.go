package main

import (
	"errors"
	"sort"
	"strconv"
	"time"

	"bridge-tab/http/middleware"
	play_app "bridge-tab/internal/play/application"
	play_query "bridge-tab/internal/play/application/query"
	play_domain "bridge-tab/internal/play/domain"
	play_infra "bridge-tab/internal/play/infrastructure"
	score_registration_infra "bridge-tab/internal/score-registration/infrastructure"
	tournament_management_query "bridge-tab/internal/tournament-management/application/query"
	tournament_domain "bridge-tab/internal/tournament-management/domain"
	tournament_management_infra "bridge-tab/internal/tournament-management/infrastructure"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

const holdRoundCookie = "holdRound"

func GetTeamPlay(c *fiber.Ctx) error {
	id, err := gameSessionID(c)
	if err != nil {
		return notFound(c, "Niepoprawny identyfikator rozgrywki")
	}
	loaded, err := loadTeamMatches(c, id)
	if err != nil {
		log.Debug(err)
		return notFound(c, translateError(err))
	}
	hold := readHoldRound(c)
	if !loaded.tournament.ScoreEntryMode.ContestantsCanRegister() {
		hold = 0
	}
	play := labeledPlay(loaded, hold)
	if isHTMX(c) {
		if play.Finished {
			c.Set("HX-Redirect", path("/round-registration/"+id))
			return c.SendStatus(fiber.StatusNoContent)
		}
		if loaded.tournament.ScoreEntryMode.ContestantsCanRegister() {
			return c.SendStatus(fiber.StatusNoContent)
		}
		return renderTeamPlay(c, id, loaded, play, false)
	}
	return renderTeamPlay(c, id, loaded, play, true)
}

func PostTeamAdvance(c *fiber.Ctx) error {
	id, err := gameSessionID(c)
	if err != nil {
		return notFound(c, "Niepoprawny identyfikator rozgrywki")
	}
	clearHoldRound(c)
	return c.Redirect(path("/round-registration/"+id), fiber.StatusSeeOther)
}

func GetTeamBoard(c *fiber.Ctx) error {
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
	play, ok := teamPlayForDeal(loaded.teamNumber, loaded.tournament.FinishedAt != "", loaded.matches, dealNo)
	if !ok {
		return notFound(c, "To rozdanie nie jest grane przez tę drużynę")
	}
	play = labelTeamPlay(play, loaded.teams)
	board, ok := findBoard(play.Boards, dealNo)
	if !ok {
		return notFound(c, "To rozdanie nie jest grane przez tę drużynę")
	}
	if !loaded.tournament.ScoreEntryMode.ContestantsCanReadScores() {
		return notFound(c, scoreEntryDenied(loaded.tournament.ScoreEntryMode))
	}
	data := boardPageData(id, loaded.tournament, play, board, c.Locals("user").(middleware.UserMetadata).Id)
	if board.Recorded {
		return c.Render("team-board-review", data, "layout")
	}
	return c.Render("team-board", data, "layout")
}

func PostTeamVerify(c *fiber.Ctx) error {
	id, err := gameSessionID(c)
	if err != nil {
		return playFragmentError(c, errors.New("niepoprawny identyfikator rozgrywki"))
	}
	loaded, err := loadTeamMatches(c, id)
	if err != nil {
		log.Debug(err)
		return playFragmentError(c, err)
	}
	body := new(boardForm)
	if err := c.BodyParser(body); err != nil {
		log.Debug(err)
		return playFragmentError(c, errors.New("błąd przy przetwarzaniu formularza"))
	}
	play, board, err := boardForForm(loaded, body.DealNo)
	if err != nil {
		return playFragmentError(c, err)
	}
	if !loaded.tournament.ScoreEntryMode.ContestantsCanRegister() {
		return playFragmentError(c, errors.New(scoreEntryDenied(loaded.tournament.ScoreEntryMode)))
	}
	if board.Recorded {
		return playFragmentError(c, play_domain.ErrBoardResultAlreadyRecorded)
	}
	return c.Render("team-confirm", fiber.Map{
		"GameSessionId":     id,
		"DealNo":            body.DealNo,
		"ContractLevel":     body.ContractLevel,
		"ContractSuit":      body.ContractSuit,
		"ContractModifier":  body.ContractModifier,
		"Tricks":            body.Tricks,
		"Declarer":          body.Declarer,
		"OpeningLeadFigure": body.OpeningLeadFigure,
		"OpeningLeadSuit":   body.OpeningLeadSuit,
		"ContractHTML":      formatCallHTML(contractValue(body)),
		"LeadHTML":          formatLeadHTML(leadValue(body)),
		"NSTeamNumber":      play.NsTeamNumber,
		"EWTeamNumber":      play.EwTeamNumber,
	})
}

func PostTeamAdd(c *fiber.Ctx) error {
	id, err := gameSessionID(c)
	if err != nil {
		return playFragmentError(c, errors.New("niepoprawny identyfikator rozgrywki"))
	}
	loaded, err := loadTeamMatches(c, id)
	if err != nil {
		log.Debug(err)
		return playFragmentError(c, err)
	}
	body := new(boardForm)
	if err := c.BodyParser(body); err != nil {
		log.Debug(err)
		return playFragmentError(c, errors.New("błąd przy przetwarzaniu formularza"))
	}
	play, board, err := boardForForm(loaded, body.DealNo)
	if err != nil {
		return playFragmentError(c, err)
	}
	if !loaded.tournament.ScoreEntryMode.ContestantsCanRegister() {
		return playFragmentError(c, errors.New(scoreEntryDenied(loaded.tournament.ScoreEntryMode)))
	}
	if board.Recorded {
		return playFragmentError(c, play_domain.ErrBoardResultAlreadyRecorded)
	}
	playerId := c.Locals("user").(middleware.UserMetadata).Id
	cmd := play_app.RecordBoardResultCommand{
		PlaySessionId:    id,
		PlayerId:         playerId,
		VersusTeamNumber: play.OpponentNumber,
		DealNo:           body.DealNo,
		Contract:         contractValue(body),
		Tricks:           body.Tricks,
		Declarer:         body.Declarer,
		OpeningLead:      leadValue(body),
	}
	if err := cmd.Execute(boardResultRepo(c), playSessionRepo(c), playedResultRepo(c)); err != nil {
		log.Debug(err)
		return playFragmentError(c, err)
	}
	matches, err := listMatches(c, id)
	if err != nil {
		log.Debug(err)
		return playFragmentError(c, err)
	}
	if teamRoundComplete(matchesForRound(matches, loaded.teamNumber, play.RoundNo)) {
		setHoldRound(c, play.RoundNo)
	}
	c.Set("HX-Redirect", path("/round-registration/"+id+"/boards/"+strconv.Itoa(body.DealNo)))
	if !isHTMX(c) {
		return c.Redirect(path("/round-registration/"+id+"/boards/"+strconv.Itoa(body.DealNo)), fiber.StatusSeeOther)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func GetTeamResults(c *fiber.Ctx) error {
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
	if !loaded.tournament.ScoreEntryMode.ContestantsCanReadScores() {
		return notFound(c, scoreEntryDenied(loaded.tournament.ScoreEntryMode))
	}
	play, ok := teamPlayForDeal(loaded.teamNumber, loaded.tournament.FinishedAt != "", loaded.matches, dealNo)
	if !ok {
		return notFound(c, "To rozdanie nie jest grane przez tę drużynę")
	}
	board, ok := findBoard(play.Boards, dealNo)
	if !ok || !board.Recorded {
		return notFound(c, "Wyniki będą dostępne po wpisaniu rozdania")
	}
	others := make([]otherResult, 0)
	for _, match := range loaded.matches {
		if match.DealNo != dealNo || match.Contract == "" || match.TableNumber == play.TableNumber {
			continue
		}
		others = append(others, otherResult{
			TableNumber:  match.TableNumber,
			ContractHTML: formatCallHTML(match.Contract),
			Declarer:     match.Declarer,
			Tricks:       match.Tricks,
			LeadHTML:     formatLeadHTML(match.OpeningLead),
		})
	}
	sort.Slice(others, func(i, j int) bool {
		return others[i].TableNumber < others[j].TableNumber
	})
	return c.Render("team-results", fiber.Map{
		"DealNo":  dealNo,
		"Results": others,
	})
}

type loadedTeam struct {
	tournament *tournament_domain.TournamentDto
	teamNumber int
	matches    []play_query.ScheduledMatchDto
	teams      []tournament_domain.TeamDto
}

func loadTeamMatches(c *fiber.Ctx, gameSessionId string) (loadedTeam, error) {
	tournament, err := loadTournament(c, gameSessionId)
	if err != nil {
		return loadedTeam{}, err
	}
	playerId := c.Locals("user").(middleware.UserMetadata).Id
	teamNumber, err := (&play_query.GetTeamNumberByMemberQuery{
		PlaySessionId: gameSessionId,
		PlayerId:      playerId,
	}).Execute(playSessionRepo(c))
	if err != nil {
		return loadedTeam{}, err
	}
	matches, err := listMatches(c, gameSessionId)
	if err != nil {
		return loadedTeam{}, err
	}
	teams, err := (&tournament_management_query.ListTeamsQuery{TournamentId: gameSessionId}).Execute(&tournament_management_infra.PostgresTeamReadRepository{
		Ctx: c.UserContext(),
		Tx:  middleware.GetTransaction(c),
	})
	if err != nil {
		return loadedTeam{}, err
	}
	return loadedTeam{
		tournament: tournament,
		teamNumber: teamNumber,
		matches:    matches,
		teams:      teams,
	}, nil
}

func labeledPlay(loaded loadedTeam, holdRound int) teamPlayView {
	play := buildTeamPlay(loaded.teamNumber, loaded.tournament.FinishedAt != "", loaded.matches, holdRound)
	return labelTeamPlay(play, loaded.teams)
}

func boardForForm(loaded loadedTeam, dealNo int) (teamPlayView, boardView, error) {
	play, ok := teamPlayForDeal(loaded.teamNumber, loaded.tournament.FinishedAt != "", loaded.matches, dealNo)
	if !ok {
		return teamPlayView{}, boardView{}, errors.New("to rozdanie nie jest grane przez tę drużynę")
	}
	board, ok := findBoard(play.Boards, dealNo)
	if !ok {
		return teamPlayView{}, boardView{}, errors.New("to rozdanie nie jest grane przez tę drużynę")
	}
	return play, board, nil
}

func matchesForRound(matches []play_query.ScheduledMatchDto, teamNumber, roundNo int) []play_query.ScheduledMatchDto {
	round := make([]play_query.ScheduledMatchDto, 0)
	for _, match := range matches {
		if match.RoundNo != roundNo {
			continue
		}
		if match.NsTeamNumber != teamNumber && match.EwTeamNumber != teamNumber {
			continue
		}
		round = append(round, match)
	}
	return round
}

func renderTeamPlay(c *fiber.Ctx, id string, loaded loadedTeam, play teamPlayView, withLayout bool) error {
	canCall := false
	if loaded.tournament.StartedAt != "" {
		var err error
		canCall, err = umpireAvailable(c, id)
		if err != nil {
			log.Debug(err)
			return notFound(c, "Błąd przy pobieraniu wezwania")
		}
	}
	data := fiber.Map{
		"Title":             loaded.tournament.Name,
		"GameSessionId":     id,
		"TournamentName":    loaded.tournament.Name,
		"CanCallUmpire":     canCall,
		"UserId":            c.Locals("user").(middleware.UserMetadata).Id,
		"TeamLabel":         play.TeamLabel,
		"RoundNo":           play.RoundNo,
		"TableNumber":       play.TableNumber,
		"Direction":         play.Direction,
		"OpponentLabel":     play.OpponentLabel,
		"Boards":            play.Boards,
		"Finished":          play.Finished,
		"AllDone":           play.AllDone,
		"Empty":             play.Empty,
		"RoundComplete":     play.RoundComplete,
		"HasNextRound":      play.HasNextRound,
		"NextRoundNo":       play.NextRoundNo,
		"NextTableNumber":   play.NextTableNumber,
		"NextDirection":     play.NextDirection,
		"NextOpponentLabel": play.NextOpponentLabel,
		"Held":              play.Held,
		"CanRegister":       loaded.tournament.ScoreEntryMode.ContestantsCanRegister(),
		"CanReadScores":     loaded.tournament.ScoreEntryMode.ContestantsCanReadScores(),
	}
	if !withLayout {
		return c.Render("team-play-round", data)
	}
	return c.Render("team-play", data, "layout")
}

func scoreEntryDenied(mode tournament_domain.ScoreEntryMode) string {
	switch mode {
	case tournament_domain.ScoreEntryUmpire:
		return "Wyniki wpisuje sędzia"
	case tournament_domain.ScoreEntryTablet:
		return "Wyniki wpisuje tablet"
	case tournament_domain.ScoreEntryContestant:
		return "Wyniki wpisują zawodnicy"
	default:
		return "Wyniki nie są dostępne"
	}
}

func boardPageData(id string, tournament *tournament_domain.TournamentDto, play teamPlayView, board boardView, userId string) fiber.Map {
	return fiber.Map{
		"Title":          "Rozdanie " + strconv.Itoa(board.DealNo),
		"GameSessionId":  id,
		"TournamentName": tournament.Name,
		"UserId":         userId,
		"TeamLabel":      play.TeamLabel,
		"TableNumber":    play.TableNumber,
		"RoundNo":        play.RoundNo,
		"Direction":      play.Direction,
		"OpponentLabel":  play.OpponentLabel,
		"DealNo":         board.DealNo,
		"ContractHTML":   board.ContractHTML,
		"Tricks":         board.Tricks,
		"Declarer":       board.Declarer,
		"LeadHTML":       board.LeadHTML,
	}
}

func gameSessionID(c *fiber.Ctx) (string, error) {
	id := c.Params("gameSessionId")
	if err := uuid.Validate(id); err != nil {
		return "", err
	}
	return id, nil
}

func dealNumber(c *fiber.Ctx) (int, error) {
	dealNo, err := strconv.Atoi(c.Params("dealNo"))
	if err != nil || dealNo < 1 {
		return 0, errors.New("nie znaleziono rozdania")
	}
	return dealNo, nil
}

func loadTournament(c *fiber.Ctx, id string) (*tournament_domain.TournamentDto, error) {
	tournament, err := (&tournament_management_query.GetTournamentById{Id: id}).Execute(&tournament_management_infra.PostgresTournamentReadRepository{
		Ctx: c.UserContext(),
		Tx:  middleware.GetTransaction(c),
	})
	if err != nil {
		return nil, err
	}
	if tournament == nil {
		return nil, errors.New("nie znaleziono turnieju")
	}
	return tournament, nil
}

func listMatches(c *fiber.Ctx, id string) ([]play_query.ScheduledMatchDto, error) {
	return (&play_query.ListMatchesQuery{PlaySessionId: id}).Execute(playSessionRepo(c), boardResultReadRepo(c))
}

func playSessionRepo(c *fiber.Ctx) *play_infra.PostgresPlaySessionRepository {
	return &play_infra.PostgresPlaySessionRepository{
		Ctx: c.UserContext(),
		Tx:  middleware.GetTransaction(c),
	}
}

func boardResultReadRepo(c *fiber.Ctx) *play_infra.PostgresBoardResultReadRepository {
	return &play_infra.PostgresBoardResultReadRepository{
		Ctx: c.UserContext(),
		Tx:  middleware.GetTransaction(c),
	}
}

func boardResultRepo(c *fiber.Ctx) *play_infra.PostgresBoardResultRepository {
	return &play_infra.PostgresBoardResultRepository{
		Ctx: c.UserContext(),
		Tx:  middleware.GetTransaction(c),
	}
}

func playedResultRepo(c *fiber.Ctx) *score_registration_infra.PostgresPlayedResultRepository {
	return &score_registration_infra.PostgresPlayedResultRepository{
		Ctx: c.UserContext(),
		Tx:  middleware.GetTransaction(c),
	}
}

func playFragmentError(c *fiber.Ctx, err error) error {
	return c.Render("error", fiber.Map{"Error": translateError(err)})
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
		Path:     cookiePath(),
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
	})
}

func clearHoldRound(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     holdRoundCookie,
		Value:    "",
		Path:     cookiePath(),
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HTTPOnly: true,
	})
}
