package main

import (
	"errors"
	"html/template"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	play_app "bridge-tab/internal/play/application"
	play_query "bridge-tab/internal/play/application/query"
	tournament_domain "bridge-tab/internal/tournament-management/domain"

	"github.com/gofiber/fiber/v2"
)

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

type tablePlayView struct {
	TournamentName string
	TableNumber    int
	RoundNo        int
	NsTeamNumber   int
	EwTeamNumber   int
	NsTeamLabel    string
	EwTeamLabel    string
	Boards         []boardView
	Finished       bool
	AllDone        bool
	Empty          bool
	RoundComplete  bool
	HasNextRound   bool
	NextRoundNo    int
	Held           bool
}

func getPlay(c *fiber.Ctx) error {
	viewData, err := loadPlay(c)
	if err != nil {
		if redirect, ok := err.(redirectError); ok {
			if isHTMX(c) {
				c.Set("HX-Redirect", path(redirect.To))
				return c.SendStatus(fiber.StatusNoContent)
			}
			return c.Redirect(path(redirect.To), fiber.StatusSeeOther)
		}
		return err
	}
	if isHTMX(c) {
		if finished, _ := viewData["Finished"].(bool); finished {
			c.Set("HX-Redirect", path("/play"))
			return c.SendStatus(fiber.StatusNoContent)
		}
		if canRegister, _ := viewData["CanRegister"].(bool); canRegister {
			return c.SendStatus(fiber.StatusNoContent)
		}
		viewData["Fragment"] = true
		return fragment(c, "play", viewData)
	}
	return playView(c, "play", viewData)
}

func postAdvance(c *fiber.Ctx) error {
	clearHoldRound(c)
	return c.Redirect(path("/play"), fiber.StatusSeeOther)
}

func getBoard(c *fiber.Ctx) error {
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
	mode := tournament.ScoreEntryMode
	if !mode.TabletCanRegister() {
		return c.Redirect(path("/play"), fiber.StatusSeeOther)
	}
	data := boardPageData(tournament, a, play, board)
	if board.Recorded {
		return playView(c, "board-review", data)
	}
	return playView(c, "board", data)
}

func postVerify(c *fiber.Ctx) error {
	play, tournament, a, err := loadTablePlayForDeal(c)
	if err != nil {
		if redirect, ok := err.(redirectError); ok {
			c.Set("HX-Redirect", path(redirect.To))
			return c.SendStatus(fiber.StatusNoContent)
		}
		return err
	}
	body := new(boardForm)
	if err := c.BodyParser(body); err != nil {
		return rejectFragment(c, "error", fiber.Map{}, errors.New("błąd przy przetwarzaniu formularza"))
	}
	board, ok := findBoard(play.Boards, body.DealNo)
	if !ok {
		return rejectFragment(c, "error", fiber.Map{}, errors.New("to rozdanie nie jest grane w tej rundzie"))
	}
	if !tournament.ScoreEntryMode.TabletCanRegister() {
		return rejectFragment(c, "error", fiber.Map{}, errors.New(tabletScoreEntryDenied(tournament.ScoreEntryMode)))
	}
	if board.Recorded {
		return rejectFragment(c, "error", fiber.Map{}, errors.New("rozdanie już zostało wpisane"))
	}

	s := newStores(c)
	match, err := (&play_query.GetMatchByTeamsQuery{
		PlaySessionId:    a.TournamentId,
		FirstTeamNumber:  play.NsTeamNumber,
		SecondTeamNumber: play.EwTeamNumber,
		DealNo:           body.DealNo,
	}).Execute(s.playSession, s.boardResults)
	if err != nil {
		return rejectFragment(c, "error", fiber.Map{}, err)
	}

	return fragment(c, "confirm-dialog", fiber.Map{
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
		"NSTeamNumber":      match.NsTeamNumber,
		"EWTeamNumber":      match.EwTeamNumber,
	})
}

func postAdd(c *fiber.Ctx) error {
	play, tournament, a, err := loadTablePlayForDeal(c)
	if err != nil {
		if redirect, ok := err.(redirectError); ok {
			c.Set("HX-Redirect", path(redirect.To))
			return c.SendStatus(fiber.StatusNoContent)
		}
		return err
	}
	body := new(boardForm)
	if err := c.BodyParser(body); err != nil {
		return rejectFragment(c, "error", fiber.Map{}, errors.New("błąd przy przetwarzaniu formularza"))
	}
	board, ok := findBoard(play.Boards, body.DealNo)
	if !ok {
		return rejectFragment(c, "error", fiber.Map{}, errors.New("to rozdanie nie jest grane w tej rundzie"))
	}
	if !tournament.ScoreEntryMode.TabletCanRegister() {
		return rejectFragment(c, "error", fiber.Map{}, errors.New(tabletScoreEntryDenied(tournament.ScoreEntryMode)))
	}
	if board.Recorded {
		return rejectFragment(c, "error", fiber.Map{}, errors.New("rozdanie już zostało wpisane"))
	}

	s := newStores(c)
	cmd := play_app.RecordBoardResultByTeamsCommand{
		PlaySessionId:    a.TournamentId,
		FirstTeamNumber:  play.NsTeamNumber,
		SecondTeamNumber: play.EwTeamNumber,
		DealNo:           body.DealNo,
		Contract:         contractValue(body),
		Tricks:           body.Tricks,
		Declarer:         body.Declarer,
		OpeningLead:      leadValue(body),
	}
	if err := cmd.Execute(s.boardResult, s.playSession, s.played); err != nil {
		return rejectFragment(c, "error", fiber.Map{}, err)
	}

	matches, err := (&play_query.ListMatchesQuery{PlaySessionId: a.TournamentId}).Execute(s.playSession, s.boardResults)
	if err != nil {
		return rejectFragment(c, "error", fiber.Map{}, err)
	}
	if roundCompleteFor(matches, a.TableNumber, play.RoundNo) {
		setHoldRound(c, play.RoundNo)
	}

	c.Set("HX-Redirect", path("/play/boards/"+strconv.Itoa(body.DealNo)))
	return c.SendStatus(fiber.StatusNoContent)
}

func getResults(c *fiber.Ctx) error {
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
	if !ok || !board.Recorded {
		return fiber.NewError(fiber.StatusNotFound, "wyniki będą dostępne po wpisaniu rozdania")
	}
	if !tournament.ScoreEntryMode.TabletCanRegister() {
		return rejectFragment(c, "error", fiber.Map{}, errors.New(tabletScoreEntryDenied(tournament.ScoreEntryMode)))
	}
	s := newStores(c)
	matches, err := (&play_query.ListMatchesQuery{PlaySessionId: a.TournamentId}).Execute(s.playSession, s.boardResults)
	if err != nil {
		return err
	}
	others := make([]otherResult, 0)
	for _, match := range matches {
		if match.DealNo != dealNo || match.Contract == "" || match.TableNumber == a.TableNumber {
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
	return fragment(c, "results", fiber.Map{
		"DealNo":  dealNo,
		"Results": others,
	})
}

func loadPlay(c *fiber.Ctx) (fiber.Map, error) {
	play, tournament, a, err := loadTablePlay(c)
	if err != nil {
		return nil, err
	}
	return playMap(tournament, play, a), nil
}

func playMap(tournament *tournament_domain.TournamentDto, play tablePlayView, a assignment) fiber.Map {
	mode := tournament.ScoreEntryMode
	held := play.Held
	if !mode.TabletCanRegister() && play.AllDone {
		held = true
	}
	return fiber.Map{
		"Title":          "Stół " + strconv.Itoa(a.TableNumber),
		"TournamentName": tournament.Name,
		"TableNumber":    play.TableNumber,
		"RoundNo":        play.RoundNo,
		"NsTeamNumber":   play.NsTeamNumber,
		"EwTeamNumber":   play.EwTeamNumber,
		"NsTeamLabel":    play.NsTeamLabel,
		"EwTeamLabel":    play.EwTeamLabel,
		"Boards":         play.Boards,
		"Finished":       play.Finished,
		"AllDone":        play.AllDone,
		"Empty":          play.Empty,
		"RoundComplete":  play.RoundComplete,
		"HasNextRound":   play.HasNextRound,
		"NextRoundNo":    play.NextRoundNo,
		"Held":           held,
		"PlayScreen":     true,
		"CanRegister":    mode.TabletCanRegister(),
		"CanReadScores":  mode.TabletCanReadScores(),
	}
}

func tabletScoreEntryDenied(mode tournament_domain.ScoreEntryMode) string {
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

func loadTablePlay(c *fiber.Ctx) (tablePlayView, *tournament_domain.TournamentDto, assignment, error) {
	return loadTablePlayWithHold(c, readHoldRound(c))
}

func loadTablePlayForDeal(c *fiber.Ctx) (tablePlayView, *tournament_domain.TournamentDto, assignment, error) {
	dealNo, _ := strconv.Atoi(c.Params("dealNo"))
	if dealNo == 0 {
		dealNo, _ = strconv.Atoi(c.FormValue("dealNo"))
	}
	a, ok := readAssignment(c)
	if !ok {
		return tablePlayView{}, nil, assignment{}, redirectError{To: "/setup"}
	}
	s := newStores(c)
	tournament, err := findTournament(s, a.TournamentId)
	if err != nil {
		return tablePlayView{}, nil, a, err
	}
	if tournament == nil {
		return tablePlayView{}, nil, a, redirectError{To: "/setup"}
	}
	if tournament.StartedAt == "" {
		return tablePlayView{}, tournament, a, redirectError{To: "/"}
	}
	matches, err := (&play_query.ListMatchesQuery{PlaySessionId: a.TournamentId}).Execute(s.playSession, s.boardResults)
	if err != nil {
		return tablePlayView{}, tournament, a, err
	}
	if dealNo > 0 {
		if play, ok := playForDeal(a.TableNumber, tournament.FinishedAt != "", matches, dealNo); ok {
			return labeledPlay(s, a.TournamentId, play, tournament, a)
		}
	}
	hold := readHoldRound(c)
	if !tournament.ScoreEntryMode.TabletCanRegister() {
		hold = 0
	}
	play := buildTablePlay(a.TableNumber, tournament.FinishedAt != "", matches, hold)
	return labeledPlay(s, a.TournamentId, play, tournament, a)
}

func loadTablePlayWithHold(c *fiber.Ctx, holdRound int) (tablePlayView, *tournament_domain.TournamentDto, assignment, error) {
	a, ok := readAssignment(c)
	if !ok {
		return tablePlayView{}, nil, assignment{}, redirectError{To: "/setup"}
	}
	s := newStores(c)
	tournament, err := findTournament(s, a.TournamentId)
	if err != nil {
		return tablePlayView{}, nil, a, err
	}
	if tournament == nil {
		return tablePlayView{}, nil, a, redirectError{To: "/setup"}
	}
	if tournament.StartedAt == "" {
		return tablePlayView{}, tournament, a, redirectError{To: "/"}
	}
	matches, err := (&play_query.ListMatchesQuery{PlaySessionId: a.TournamentId}).Execute(s.playSession, s.boardResults)
	if err != nil {
		return tablePlayView{}, tournament, a, err
	}
	if !tournament.ScoreEntryMode.TabletCanRegister() {
		holdRound = 0
	}
	play := buildTablePlay(a.TableNumber, tournament.FinishedAt != "", matches, holdRound)
	return labeledPlay(s, a.TournamentId, play, tournament, a)
}

func labeledPlay(s stores, tournamentId string, play tablePlayView, tournament *tournament_domain.TournamentDto, a assignment) (tablePlayView, *tournament_domain.TournamentDto, assignment, error) {
	teams, err := s.teams.FindAll(&tournamentId)
	if err != nil {
		return play, tournament, a, err
	}
	play.NsTeamLabel = teamLabel(teams, play.NsTeamNumber)
	play.EwTeamLabel = teamLabel(teams, play.EwTeamNumber)
	return play, tournament, a, nil
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

func buildTablePlay(tableNumber int, finished bool, matches []play_query.ScheduledMatchDto, holdRound int) tablePlayView {
	byRound, roundNos := roundsForTable(matches, tableNumber)
	view := tablePlayView{
		TableNumber: tableNumber,
		Finished:    finished,
		Empty:       len(roundNos) == 0,
	}
	if view.Empty {
		return view
	}

	firstIncomplete := 0
	allDone := true
	for _, roundNo := range roundNos {
		if !roundComplete(byRound[roundNo]) {
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
		if held, ok := byRound[holdRound]; ok && roundComplete(held) {
			chosen = holdRound
			view.Held = true
		}
	}

	view.AllDone = allDone
	applyRound(&view, byRound, roundNos, chosen)
	return view
}

func playForDeal(tableNumber int, finished bool, matches []play_query.ScheduledMatchDto, dealNo int) (tablePlayView, bool) {
	byRound, roundNos := roundsForTable(matches, tableNumber)
	for _, roundNo := range roundNos {
		for _, match := range byRound[roundNo] {
			if match.DealNo == dealNo {
				view := tablePlayView{TableNumber: tableNumber, Finished: finished}
				applyRound(&view, byRound, roundNos, roundNo)
				return view, true
			}
		}
	}
	return tablePlayView{}, false
}

func applyRound(view *tablePlayView, byRound map[int][]play_query.ScheduledMatchDto, roundNos []int, chosen int) {
	boards := byRound[chosen]
	sort.Slice(boards, func(i, j int) bool {
		return boards[i].DealNo < boards[j].DealNo
	})
	view.RoundNo = chosen
	view.NsTeamNumber = boards[0].NsTeamNumber
	view.EwTeamNumber = boards[0].EwTeamNumber
	view.Boards = toBoardViews(boards)
	view.RoundComplete = roundComplete(boards)
	view.HasNextRound = false
	view.NextRoundNo = 0
	for _, roundNo := range roundNos {
		if roundNo > chosen {
			view.HasNextRound = true
			view.NextRoundNo = roundNo
			break
		}
	}
}

func roundsForTable(matches []play_query.ScheduledMatchDto, tableNumber int) (map[int][]play_query.ScheduledMatchDto, []int) {
	byRound := map[int][]play_query.ScheduledMatchDto{}
	roundNos := make([]int, 0)
	for _, match := range matches {
		if match.TableNumber != tableNumber {
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

func roundComplete(matches []play_query.ScheduledMatchDto) bool {
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

func roundCompleteFor(matches []play_query.ScheduledMatchDto, tableNumber, roundNo int) bool {
	var round []play_query.ScheduledMatchDto
	for _, match := range matches {
		if match.TableNumber == tableNumber && match.RoundNo == roundNo {
			round = append(round, match)
		}
	}
	return roundComplete(round)
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

func boardPageData(tournament *tournament_domain.TournamentDto, a assignment, play tablePlayView, board boardView) fiber.Map {
	return fiber.Map{
		"Title":             "Rozdanie " + strconv.Itoa(board.DealNo),
		"TournamentName":    tournament.Name,
		"TableNumber":       a.TableNumber,
		"RoundNo":           play.RoundNo,
		"NsTeamNumber":      play.NsTeamNumber,
		"EwTeamNumber":      play.EwTeamNumber,
		"NsTeamLabel":       play.NsTeamLabel,
		"EwTeamLabel":       play.EwTeamLabel,
		"DealNo":            board.DealNo,
		"Recorded":          board.Recorded,
		"Contract":          board.Contract,
		"ContractHTML":      board.ContractHTML,
		"Tricks":            board.Tricks,
		"Declarer":          board.Declarer,
		"OpeningLead":       board.OpeningLead,
		"LeadHTML":          board.LeadHTML,
		"PlayScreen":        true,
		"ShowVulnerability": true,
		"NsVulnerable":      board.NsVulnerable,
		"EwVulnerable":      board.EwVulnerable,
	}
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

type redirectError struct {
	To string
}

func (e redirectError) Error() string {
	return "redirect " + e.To
}
