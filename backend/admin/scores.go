package main

import (
	"sort"
	"strconv"
	"strings"

	"bridge-tab/internal/idutil"
	rounds_query "bridge-tab/internal/rounds-registration/application/query"
	rounds_domain "bridge-tab/internal/rounds-registration/domain"
	score_app "bridge-tab/internal/score-registration/application"
	score_query "bridge-tab/internal/score-registration/application/query"
	score_domain "bridge-tab/internal/score-registration/domain"
	tournament_query "bridge-tab/internal/tournament-management/application/query"
	tournament_domain "bridge-tab/internal/tournament-management/domain"

	"github.com/gofiber/fiber/v2"
)

type standingRow struct {
	Rank       int
	Number     int
	Name       string
	Total      string
	Percentage string
}

type boardLine struct {
	NS          string
	EW          string
	Contract    string
	Declarer    string
	Tricks      string
	OpeningLead string
	Points      int
	NsAward     string
	EwAward     string
}

type dealTable struct {
	Deal       int
	Vulnerable string
	Lines      []boardLine
}

type scoreGroup struct {
	Label     string
	Deals     []dealTable
	PageBreak bool
}

type sheetTeam struct {
	Number int
	Name   string
}

type sheetRow struct {
	Deal    int
	Set     string
	SetSpan int
	Cells   []string
}

type pairSheet struct {
	Teams  []sheetTeam
	Rows   []sheetRow
	Totals []string
}

type setSheetRow struct {
	Label  string
	Boards string
	Cells  []string
}

type setSheet struct {
	Teams  []sheetTeam
	Rows   []setSheetRow
	Totals []string
}

type roundKey struct {
	deal int
	ns   int
	ew   int
}

func getScores(c *fiber.Ctx) error {
	s := newStores(c)
	tournament, err := loadTournament(c, s)
	if err != nil {
		return err
	}
	data, err := scoresPage(s, tournament)
	if err != nil {
		return err
	}
	return view(c, "scores", data)
}

func postCalculate(c *fiber.Ctx) error {
	s := newStores(c)
	tournament, err := loadTournament(c, s)
	if err != nil {
		return err
	}
	data, err := scoresPage(s, tournament)
	if err != nil {
		return err
	}
	cmd := score_app.CalculateStandingsCommand{
		TournamentId: tournament.Id,
		Method:       score_domain.MethodMatchpoints,
	}
	if err := cmd.Execute(s.played, s.scores); err != nil {
		return reject(c, "scores", data, err)
	}
	return c.Redirect("/tournaments/"+tournament.Id+"/scores", fiber.StatusSeeOther)
}

func postScoresFinish(c *fiber.Ctx) error {
	return finishTournament(c, "/tournaments/"+c.Params("id")+"/scores")
}

func scoresPage(s stores, tournament *tournament_domain.TournamentDto) (fiber.Map, error) {
	standings, err := (&score_query.GetStandingsQuery{
		TournamentId: tournament.Id,
		Method:       score_domain.MethodMatchpoints,
	}).Execute(s.scoreRead)
	if err != nil {
		return nil, err
	}
	boards, err := (&score_query.ListBoardScoresQuery{
		TournamentId: tournament.Id,
		Method:       score_domain.MethodMatchpoints,
	}).Execute(s.scoreRead)
	if err != nil {
		return nil, err
	}
	sets, err := (&tournament_query.ListSetsQuery{TournamentId: tournament.Id}).Execute(s.sets)
	if err != nil {
		return nil, err
	}
	teams, err := (&tournament_query.ListTeamsQuery{TournamentId: tournament.Id}).Execute(s.teams)
	if err != nil {
		return nil, err
	}
	rounds, err := (&rounds_query.ListRoundsQuery{GameSessionId: tournament.Id}).Execute(s.sessionRead)
	if err != nil {
		return nil, err
	}
	protocols, err := (&tournament_query.ListBoardProtocolsQuery{TournamentId: tournament.Id}).Execute(s.boards)
	if err != nil {
		return nil, err
	}

	rows := make([]standingRow, len(standings))
	for i, standing := range standings {
		rows[i] = standingRow{
			Rank:       standing.Rank,
			Number:     standing.TeamNumber,
			Name:       standing.TeamName,
			Total:      formatScore(standing.TotalScore),
			Percentage: formatPct(standing.Percentage),
		}
	}

	pairs, setsSheet := resultSheets(sets, boards, teams)

	data := pageData(tournament)
	data["Standings"] = rows
	data["Pairs"] = pairs
	data["SetsSheet"] = setsSheet
	data["Groups"] = scoreGroups(sets, boards, teams, rounds, protocols)
	data["HasScores"] = len(rows) > 0 || len(boards) > 0
	return data, nil
}

func scoreGroups(sets []tournament_domain.SetDto, boards []score_domain.BoardScoreDto, teams []tournament_domain.TeamDto, rounds []rounds_domain.PlayedRoundDto, protocols []tournament_domain.BoardProtocolDto) []scoreGroup {
	if len(boards) == 0 {
		return nil
	}
	played := indexRounds(rounds)
	vulnerable := indexVulnerable(protocols)
	if len(sets) == 0 {
		return []scoreGroup{boardGroup("Rozdania", boards, teams, played, vulnerable, false)}
	}

	assigned := map[int]bool{}
	for _, set := range sets {
		for _, boardNo := range set.BoardNos {
			assigned[boardNo] = true
		}
	}

	groups := make([]scoreGroup, 0, len(sets)+1)
	for _, set := range sets {
		var lines []score_domain.BoardScoreDto
		for _, board := range boards {
			if containsInt(set.BoardNos, board.DealNo) {
				lines = append(lines, board)
			}
		}
		if len(lines) == 0 {
			continue
		}
		groups = append(groups, boardGroup(set.Label, lines, teams, played, vulnerable, true))
	}

	var rest []score_domain.BoardScoreDto
	for _, board := range boards {
		if !assigned[board.DealNo] {
			rest = append(rest, board)
		}
	}
	if len(rest) > 0 {
		groups = append(groups, boardGroup("Poza setami", rest, teams, played, vulnerable, true))
	}
	return groups
}

func boardGroup(label string, boards []score_domain.BoardScoreDto, teams []tournament_domain.TeamDto, played map[roundKey]rounds_domain.PlayedRoundDto, vulnerable map[int]string, pageBreak bool) scoreGroup {
	sorted := append([]score_domain.BoardScoreDto(nil), boards...)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].DealNo == sorted[j].DealNo {
			return sorted[i].NsTeamNumber < sorted[j].NsTeamNumber
		}
		return sorted[i].DealNo < sorted[j].DealNo
	})
	deals := make([]dealTable, 0)
	for _, board := range sorted {
		line := boardLine{
			NS:      teamCaption(teams, board.NsTeamId, board.NsTeamNumber),
			EW:      teamCaption(teams, board.EwTeamId, board.EwTeamNumber),
			Points:  board.NsPoints,
			NsAward: formatScore(board.NsAward),
			EwAward: formatScore(board.EwAward),
		}
		if round, ok := played[roundKey{board.DealNo, board.NsTeamNumber, board.EwTeamNumber}]; ok {
			line.Contract = displayContract(round.Contract)
			line.Declarer = round.Declarer
			line.OpeningLead = round.OpeningLead
			if round.Contract != "" && round.Contract != "Pass" {
				line.Tricks = strconv.Itoa(round.Tricks)
			}
		}
		if len(deals) == 0 || deals[len(deals)-1].Deal != board.DealNo {
			deals = append(deals, dealTable{
				Deal:       board.DealNo,
				Vulnerable: vulnerable[board.DealNo],
				Lines:      []boardLine{line},
			})
			continue
		}
		last := &deals[len(deals)-1]
		last.Lines = append(last.Lines, line)
	}
	return scoreGroup{
		Label:     label,
		Deals:     deals,
		PageBreak: pageBreak,
	}
}

type teamAwardKey struct {
	team string
	deal int
}

func resultSheets(sets []tournament_domain.SetDto, boards []score_domain.BoardScoreDto, teams []tournament_domain.TeamDto) (pairSheet, setSheet) {
	if len(boards) == 0 {
		return pairSheet{}, setSheet{}
	}
	columns := sheetTeams(boards, teams)
	awards := map[teamAwardKey]float64{}
	for _, board := range boards {
		awards[teamAwardKey{canonId(board.NsTeamId), board.DealNo}] = board.NsAward
		awards[teamAwardKey{canonId(board.EwTeamId), board.DealNo}] = board.EwAward
	}

	rows := make([]sheetRow, 0)
	for _, placed := range orderedDeals(sets, boards) {
		cells := make([]string, len(columns))
		for i, team := range columns {
			award, ok := awards[teamAwardKey{team.id, placed.deal}]
			if !ok {
				continue
			}
			cells[i] = formatScore(award)
		}
		rows = append(rows, sheetRow{
			Deal:  placed.deal,
			Set:   placed.set,
			Cells: cells,
		})
	}
	markSetSpans(rows)

	totals := make([]float64, len(columns))
	for _, board := range boards {
		for i, team := range columns {
			if team.id == canonId(board.NsTeamId) {
				totals[i] += board.NsAward
			}
			if team.id == canonId(board.EwTeamId) {
				totals[i] += board.EwAward
			}
		}
	}
	totalCells := make([]string, len(columns))
	for i, total := range totals {
		totalCells[i] = formatScore(total)
	}

	return pairSheet{
		Teams:  sheetTeamView(columns),
		Rows:   rows,
		Totals: totalCells,
	}, setTotalsSheet(sets, columns, awards)
}

type placedDeal struct {
	deal int
	set  string
}

func orderedDeals(sets []tournament_domain.SetDto, boards []score_domain.BoardScoreDto) []placedDeal {
	scored := map[int]bool{}
	for _, board := range boards {
		scored[board.DealNo] = true
	}
	seen := map[int]bool{}
	out := make([]placedDeal, 0)
	for _, set := range sets {
		nums := append([]int(nil), set.BoardNos...)
		sort.Ints(nums)
		for _, deal := range nums {
			if !scored[deal] || seen[deal] {
				continue
			}
			seen[deal] = true
			out = append(out, placedDeal{deal: deal, set: set.Label})
		}
	}
	var rest []int
	for deal := range scored {
		if !seen[deal] {
			rest = append(rest, deal)
		}
	}
	sort.Ints(rest)
	label := ""
	if len(sets) > 0 {
		label = "poza"
	}
	for _, deal := range rest {
		out = append(out, placedDeal{deal: deal, set: label})
	}
	return out
}

func markSetSpans(rows []sheetRow) {
	for i := 0; i < len(rows); {
		j := i + 1
		for j < len(rows) && rows[j].Set == rows[i].Set {
			j++
		}
		rows[i].SetSpan = j - i
		i = j
	}
}

type sheetTeamCol struct {
	id     string
	number int
	name   string
}

func sheetTeams(boards []score_domain.BoardScoreDto, teams []tournament_domain.TeamDto) []sheetTeamCol {
	seen := map[string]sheetTeamCol{}
	add := func(id string, number int) {
		id = canonId(id)
		if _, ok := seen[id]; ok {
			return
		}
		col := sheetTeamCol{id: id, number: number}
		for _, team := range teams {
			if idutil.SameId(team.Id, id) {
				col.number = team.Number
				col.name = team.Name
				break
			}
		}
		seen[id] = col
	}
	for _, board := range boards {
		add(board.NsTeamId, board.NsTeamNumber)
		add(board.EwTeamId, board.EwTeamNumber)
	}
	cols := make([]sheetTeamCol, 0, len(seen))
	for _, col := range seen {
		cols = append(cols, col)
	}
	sort.Slice(cols, func(i, j int) bool {
		if cols[i].number == cols[j].number {
			return cols[i].id < cols[j].id
		}
		return cols[i].number < cols[j].number
	})
	return cols
}

func sheetTeamView(cols []sheetTeamCol) []sheetTeam {
	out := make([]sheetTeam, len(cols))
	for i, col := range cols {
		out[i] = sheetTeam{Number: col.number, Name: col.name}
	}
	return out
}

func setTotalsSheet(sets []tournament_domain.SetDto, columns []sheetTeamCol, awards map[teamAwardKey]float64) setSheet {
	if len(sets) == 0 || len(columns) == 0 {
		return setSheet{}
	}
	rows := make([]setSheetRow, 0, len(sets))
	totals := make([]float64, len(columns))
	for _, set := range sets {
		cells := make([]string, len(columns))
		any := false
		for i, team := range columns {
			var sum float64
			played := false
			for _, deal := range set.BoardNos {
				award, ok := awards[teamAwardKey{team.id, deal}]
				if !ok {
					continue
				}
				sum += award
				played = true
			}
			if !played {
				continue
			}
			cells[i] = formatScore(sum)
			totals[i] += sum
			any = true
		}
		if !any {
			continue
		}
		rows = append(rows, setSheetRow{
			Label:  set.Label,
			Boards: formatBoardSpan(set.BoardNos),
			Cells:  cells,
		})
	}
	if len(rows) == 0 {
		return setSheet{}
	}
	totalCells := make([]string, len(columns))
	for i, total := range totals {
		totalCells[i] = formatScore(total)
	}
	return setSheet{
		Teams:  sheetTeamView(columns),
		Rows:   rows,
		Totals: totalCells,
	}
}

func canonId(id string) string {
	return strings.ToLower(id)
}

func formatBoardSpan(nums []int) string {
	if len(nums) == 0 {
		return ""
	}
	sorted := append([]int(nil), nums...)
	sort.Ints(sorted)
	if len(sorted) > 1 && sorted[len(sorted)-1]-sorted[0]+1 == len(sorted) {
		return strconv.Itoa(sorted[0]) + "–" + strconv.Itoa(sorted[len(sorted)-1])
	}
	return joinInts(sorted)
}

func indexRounds(rounds []rounds_domain.PlayedRoundDto) map[roundKey]rounds_domain.PlayedRoundDto {
	out := make(map[roundKey]rounds_domain.PlayedRoundDto, len(rounds))
	for _, round := range rounds {
		out[roundKey{round.DealNo, round.NsTeamNumber, round.EwTeamNumber}] = round
	}
	return out
}

func indexVulnerable(protocols []tournament_domain.BoardProtocolDto) map[int]string {
	out := make(map[int]string, len(protocols))
	for _, protocol := range protocols {
		out[protocol.BoardNo] = vulLabel(protocol.Vulnerable)
	}
	return out
}

func displayContract(contract string) string {
	if i := strings.Index(contract, "N"); i >= 0 && !strings.Contains(contract, "NT") {
		return contract[:i] + "NT" + contract[i+1:]
	}
	return contract
}

func containsInt(nums []int, n int) bool {
	for _, value := range nums {
		if value == n {
			return true
		}
	}
	return false
}
