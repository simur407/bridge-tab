package main

import (
	"errors"
	"net/url"
	"sort"
	"strings"
	"time"

	"bridge-tab/http/middleware"
	play_query "bridge-tab/internal/play/application/query"
	play_domain "bridge-tab/internal/play/domain"
	play_infra "bridge-tab/internal/play/infrastructure"
	umpirecall "bridge-tab/internal/umpire-call"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

type callSeat struct {
	TableNumber  int
	RoundNo      int
	NsTeamNumber int
	EwTeamNumber int
	Current      bool
}

func getUmpireCall(c *fiber.Ctx) error {
	tournamentId := c.Params("tournamentId")
	if err := uuid.Validate(tournamentId); err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	data, show, err := loadUmpireButton(c, tournamentId)
	if err != nil {
		log.Debug(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if !show {
		return c.SendString("")
	}
	return c.Render("umpire-call", data)
}

func postUmpireCall(c *fiber.Ctx) error {
	tournamentId := c.Params("tournamentId")
	if err := uuid.Validate(tournamentId); err != nil {
		return notFound(c, "Niepoprawny identyfikator turnieju")
	}
	seat, call, err := umpireContext(c, tournamentId)
	if err != nil {
		log.Debug(err)
		return notFound(c, translateError(err))
	}
	if seat.Current {
		if _, err := umpireRepo(c).Insert(umpirecall.Call{
			TournamentId: tournamentId,
			TableNumber:  seat.TableNumber,
			RoundNo:      seat.RoundNo,
			NsTeamNumber: seat.NsTeamNumber,
			EwTeamNumber: seat.EwTeamNumber,
			Source:       umpirecall.SourceContestant,
		}); err != nil {
			log.Debug(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		call, err = umpireRepo(c).LatestForTable(tournamentId, seat.TableNumber)
		if err != nil {
			log.Debug(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}
	data, show := umpireButtonView(call, tournamentId, seat)
	if isHTMX(c) {
		if !show {
			return c.SendString("")
		}
		return c.Render("umpire-call", data)
	}
	return c.Redirect(contestantReturn(c, "/tournaments/"+tournamentId+"/teams"), fiber.StatusSeeOther)
}

func umpireAvailable(c *fiber.Ctx, tournamentId string) (bool, error) {
	_, show, err := loadUmpireButton(c, tournamentId)
	return show, err
}

func loadUmpireButton(c *fiber.Ctx, tournamentId string) (fiber.Map, bool, error) {
	seat, call, err := umpireContext(c, tournamentId)
	if err != nil {
		if errors.Is(err, play_domain.ErrTeamNotFound) || errors.Is(err, play_domain.ErrPlaySessionNotProvisioned) {
			return nil, false, nil
		}
		return nil, false, err
	}
	data, show := umpireButtonView(call, tournamentId, seat)
	return data, show, nil
}

func umpireContext(c *fiber.Ctx, tournamentId string) (callSeat, *umpirecall.Call, error) {
	playerId := c.Locals("user").(middleware.UserMetadata).Id
	teamNumber, err := (&play_query.GetTeamNumberByMemberQuery{
		PlaySessionId: tournamentId,
		PlayerId:      playerId,
	}).Execute(&play_infra.PostgresPlaySessionRepository{
		Ctx: c.UserContext(),
		Tx:  middleware.GetTransaction(c),
	})
	if err != nil {
		return callSeat{}, nil, err
	}
	matches, err := (&play_query.ListMatchesQuery{PlaySessionId: tournamentId}).Execute(
		&play_infra.PostgresPlaySessionRepository{
			Ctx: c.UserContext(),
			Tx:  middleware.GetTransaction(c),
		},
		&play_infra.PostgresBoardResultReadRepository{
			Ctx: c.UserContext(),
			Tx:  middleware.GetTransaction(c),
		},
	)
	if err != nil {
		return callSeat{}, nil, err
	}
	seat, ok := seatForTeam(matches, teamNumber)
	if !ok {
		return callSeat{}, nil, nil
	}
	call, err := umpireRepo(c).LatestForTable(tournamentId, seat.TableNumber)
	if err != nil {
		return callSeat{}, nil, err
	}
	return seat, call, nil
}

func umpireButtonView(call *umpirecall.Call, tournamentId string, seat callSeat) (fiber.Map, bool) {
	state, _ := umpirecall.ButtonFrom(call, time.Now())
	if !seat.Current && state == umpirecall.ButtonIdle {
		return nil, false
	}
	if seat.TableNumber == 0 {
		return nil, false
	}
	url := "/tournaments/" + tournamentId + "/umpire-call"
	return umpireButtonData(call, url), true
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

func umpireRepo(c *fiber.Ctx) *umpirecall.Repository {
	return &umpirecall.Repository{
		Ctx: c.UserContext(),
		Tx:  middleware.GetTransaction(c),
	}
}

func seatForTeam(matches []play_query.ScheduledMatchDto, teamNumber int) (callSeat, bool) {
	if teamNumber == 0 {
		return callSeat{}, false
	}
	byRound := map[int][]play_query.ScheduledMatchDto{}
	rounds := make([]int, 0)
	for _, match := range matches {
		if match.NsTeamNumber != teamNumber && match.EwTeamNumber != teamNumber {
			continue
		}
		if _, ok := byRound[match.RoundNo]; !ok {
			rounds = append(rounds, match.RoundNo)
		}
		byRound[match.RoundNo] = append(byRound[match.RoundNo], match)
	}
	if len(rounds) == 0 {
		return callSeat{}, false
	}
	sort.Ints(rounds)
	chosen := 0
	current := false
	for _, roundNo := range rounds {
		if !roundRecorded(byRound[roundNo]) {
			chosen = roundNo
			current = true
			break
		}
	}
	if !current {
		chosen = rounds[len(rounds)-1]
	}
	board := byRound[chosen][0]
	return callSeat{
		TableNumber:  board.TableNumber,
		RoundNo:      chosen,
		NsTeamNumber: board.NsTeamNumber,
		EwTeamNumber: board.EwTeamNumber,
		Current:      current,
	}, true
}

func roundRecorded(matches []play_query.ScheduledMatchDto) bool {
	if len(matches) == 0 {
		return false
	}
	for _, match := range matches {
		if match.Contract == "" {
			return false
		}
	}
	return true
}

func contestantReturn(c *fiber.Ctx, fallback string) string {
	ref := c.Get("Referer")
	u, err := url.Parse(ref)
	if err != nil || u.Path == "" || !strings.HasPrefix(u.Path, "/") || strings.HasPrefix(u.Path, "//") {
		return path(fallback)
	}
	if u.Host != "" && !strings.EqualFold(u.Hostname(), c.Hostname()) {
		return path(fallback)
	}
	return u.RequestURI()
}

func isHTMX(c *fiber.Ctx) bool {
	return c.Get("HX-Request") == "true"
}
