package main

import (
	"html/template"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	play_query "bridge-tab/internal/play/application/query"
	tournament_domain "bridge-tab/internal/tournament-management/domain"
)

type boardView struct {
	DealNo          int
	Recorded        bool
	Contract        string
	Tricks          int
	Declarer        string
	OpeningLead     string
	ContractHTML    template.HTML
	LeadHTML        template.HTML
	Vulnerable      string
	VulnerableLabel string
	NsVulnerable    bool
	EwVulnerable    bool
}

type teamPlayView struct {
	TeamNumber         int
	TeamLabel          string
	RoundNo            int
	TableNumber        int
	Direction          string
	OpponentNumber     int
	OpponentLabel      string
	NsTeamNumber       int
	EwTeamNumber       int
	Boards             []boardView
	Finished           bool
	AllDone            bool
	Empty              bool
	RoundComplete      bool
	HasNextRound       bool
	NextRoundNo        int
	NextTableNumber    int
	NextDirection      string
	NextOpponentNumber int
	NextOpponentLabel  string
	Held               bool
}

func buildTeamPlay(teamNumber int, finished bool, matches []play_query.ScheduledMatchDto, holdRound int) teamPlayView {
	byRound, roundNos := roundsForTeam(matches, teamNumber)
	view := teamPlayView{
		TeamNumber: teamNumber,
		Finished:   finished,
		Empty:      len(roundNos) == 0,
	}
	if view.Empty {
		return view
	}

	firstIncomplete := 0
	allDone := true
	for _, roundNo := range roundNos {
		if !teamRoundComplete(byRound[roundNo]) {
			firstIncomplete = roundNo
			allDone = false
			break
		}
	}

	chosen := firstIncomplete
	if chosen == 0 {
		chosen = roundNos[len(roundNos)-1]
	}
	if holdRound > 0 {
		if held, ok := byRound[holdRound]; ok && teamRoundComplete(held) {
			chosen = holdRound
			view.Held = true
		}
	}

	view.AllDone = allDone
	applyTeamRound(&view, byRound, roundNos, chosen, teamNumber)
	return view
}

func teamPlayForDeal(teamNumber int, finished bool, matches []play_query.ScheduledMatchDto, dealNo int) (teamPlayView, bool) {
	byRound, roundNos := roundsForTeam(matches, teamNumber)
	for _, roundNo := range roundNos {
		for _, match := range byRound[roundNo] {
			if match.DealNo == dealNo {
				view := teamPlayView{TeamNumber: teamNumber, Finished: finished}
				applyTeamRound(&view, byRound, roundNos, roundNo, teamNumber)
				return view, true
			}
		}
	}
	return teamPlayView{}, false
}

func labelTeamPlay(view teamPlayView, teams []tournament_domain.TeamDto) teamPlayView {
	view.TeamLabel = teamLabel(teams, view.TeamNumber)
	view.OpponentLabel = teamLabel(teams, view.OpponentNumber)
	view.NextOpponentLabel = teamLabel(teams, view.NextOpponentNumber)
	return view
}

func applyTeamRound(view *teamPlayView, byRound map[int][]play_query.ScheduledMatchDto, roundNos []int, chosen int, teamNumber int) {
	boards := append([]play_query.ScheduledMatchDto(nil), byRound[chosen]...)
	sort.Slice(boards, func(i, j int) bool {
		return boards[i].DealNo < boards[j].DealNo
	})
	first := boards[0]
	opponent, direction := seatingOf(first, teamNumber)
	view.TeamNumber = teamNumber
	view.RoundNo = chosen
	view.TableNumber = first.TableNumber
	view.Direction = direction
	view.OpponentNumber = opponent
	view.NsTeamNumber = first.NsTeamNumber
	view.EwTeamNumber = first.EwTeamNumber
	view.Boards = toBoardViews(boards)
	view.RoundComplete = teamRoundComplete(boards)
	view.HasNextRound = false
	view.NextRoundNo = 0
	view.NextTableNumber = 0
	view.NextDirection = ""
	view.NextOpponentNumber = 0
	for _, roundNo := range roundNos {
		if roundNo <= chosen {
			continue
		}
		nextBoards := append([]play_query.ScheduledMatchDto(nil), byRound[roundNo]...)
		sort.Slice(nextBoards, func(i, j int) bool {
			return nextBoards[i].DealNo < nextBoards[j].DealNo
		})
		next := nextBoards[0]
		nextOpponent, nextDirection := seatingOf(next, teamNumber)
		view.HasNextRound = true
		view.NextRoundNo = roundNo
		view.NextTableNumber = next.TableNumber
		view.NextDirection = nextDirection
		view.NextOpponentNumber = nextOpponent
		break
	}
}

func roundsForTeam(matches []play_query.ScheduledMatchDto, teamNumber int) (map[int][]play_query.ScheduledMatchDto, []int) {
	byRound := map[int][]play_query.ScheduledMatchDto{}
	roundNos := make([]int, 0)
	if teamNumber == 0 {
		return byRound, roundNos
	}
	for _, match := range matches {
		if match.NsTeamNumber != teamNumber && match.EwTeamNumber != teamNumber {
			continue
		}
		if _, ok := byRound[match.RoundNo]; !ok {
			roundNos = append(roundNos, match.RoundNo)
		}
		byRound[match.RoundNo] = append(byRound[match.RoundNo], match)
	}
	sort.Ints(roundNos)
	return byRound, roundNos
}

func seatingOf(match play_query.ScheduledMatchDto, teamNumber int) (opponent int, direction string) {
	if match.NsTeamNumber == teamNumber {
		return match.EwTeamNumber, "NS"
	}
	return match.NsTeamNumber, "EW"
}

func teamRoundComplete(matches []play_query.ScheduledMatchDto) bool {
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

func toBoardViews(matches []play_query.ScheduledMatchDto) []boardView {
	boards := make([]boardView, len(matches))
	for i, match := range matches {
		boards[i] = boardView{
			DealNo:          match.DealNo,
			Recorded:        match.Contract != "",
			Contract:        displayContract(match.Contract),
			Tricks:          match.Tricks,
			Declarer:        match.Declarer,
			OpeningLead:     match.OpeningLead,
			ContractHTML:    formatCallHTML(match.Contract),
			LeadHTML:        formatLeadHTML(match.OpeningLead),
			Vulnerable:      match.Vulnerable,
			VulnerableLabel: vulLabel(match.Vulnerable),
			NsVulnerable:    nsVulnerable(match.Vulnerable),
			EwVulnerable:    ewVulnerable(match.Vulnerable),
		}
	}
	return boards
}

func findBoard(boards []boardView, dealNo int) (boardView, bool) {
	for _, board := range boards {
		if board.DealNo == dealNo {
			return board, true
		}
	}
	return boardView{}, false
}

func teamLabel(teams []tournament_domain.TeamDto, number int) string {
	if number == 0 {
		return ""
	}
	for _, team := range teams {
		if team.Number == number {
			name := strings.TrimSpace(team.Name)
			if name != "" {
				return strconv.Itoa(number) + " " + name
			}
			break
		}
	}
	return "Drużyna " + strconv.Itoa(number)
}

func vulLabel(value string) string {
	switch value {
	case "NS":
		return "NS"
	case "EW":
		return "EW"
	case "Both":
		return "obie"
	default:
		return "brak"
	}
}

func nsVulnerable(value string) bool {
	return value == "NS" || value == "Both"
}

func ewVulnerable(value string) bool {
	return value == "EW" || value == "Both"
}

type boardForm struct {
	DealNo            int    `form:"dealNo"`
	ContractLevel     string `form:"contractLevel"`
	ContractSuit      string `form:"contractSuit"`
	ContractModifier  string `form:"contractModifier"`
	Tricks            int    `form:"tricks"`
	Declarer          string `form:"declarer"`
	OpeningLeadSuit   string `form:"openingLeadSuit"`
	OpeningLeadFigure string `form:"openingLeadFigure"`
}

func contractValue(body *boardForm) string {
	if body.ContractModifier == "Pass" {
		return "Pass"
	}
	return body.ContractLevel + body.ContractSuit + body.ContractModifier
}

func leadValue(body *boardForm) string {
	if body.ContractModifier == "Pass" {
		return ""
	}
	return body.OpeningLeadFigure + body.OpeningLeadSuit
}

func displayContract(contract string) string {
	if contract == "" || strings.Contains(contract, "NT") {
		return contract
	}
	return strings.Replace(contract, "N", "NT", 1)
}

func formatCallHTML(contract string) template.HTML {
	if contract == "" {
		return "—"
	}
	if contract == "Pass" {
		return "Pas"
	}
	contract = displayContract(contract)
	level, suit, rest := splitContract(contract)
	symbol, class := suitSymbol(suit)
	return template.HTML(level + spanSuit(symbol, class) + rest)
}

func formatLeadHTML(lead string) template.HTML {
	if lead == "" {
		return "—"
	}
	figure, suit := splitLead(lead)
	symbol, class := suitSymbol(suit)
	return template.HTML(figure + spanSuit(symbol, class))
}

func splitContract(contract string) (level, suit, rest string) {
	if contract == "" {
		return "", "", ""
	}
	level = contract[:1]
	body := contract[1:]
	if strings.HasPrefix(body, "NT") {
		return level, "NT", body[2:]
	}
	if body == "" {
		return level, "", ""
	}
	r, size := utf8.DecodeRuneInString(body)
	return level, string(r), body[size:]
}

func splitLead(lead string) (figure, suit string) {
	if strings.HasPrefix(lead, "10") {
		return "10", lead[2:]
	}
	if lead == "" {
		return "", ""
	}
	r, size := utf8.DecodeRuneInString(lead)
	return string(r), lead[size:]
}

func suitSymbol(suit string) (string, string) {
	switch suit {
	case "C":
		return "♣", "club"
	case "D":
		return "♦", "red"
	case "H":
		return "♥", "red"
	case "S":
		return "♠", "club"
	case "N", "NT":
		return "NT", ""
	default:
		return suit, ""
	}
}

func spanSuit(symbol, class string) string {
	if symbol == "" {
		return ""
	}
	if class == "" {
		return symbol
	}
	return `<span class="suit ` + class + `">` + symbol + `</span>`
}

type otherResult struct {
	TableNumber  int
	ContractHTML template.HTML
	Declarer     string
	Tricks       int
	LeadHTML     template.HTML
}
