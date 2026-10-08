package score_registration

import (
	"errors"
	"sort"
)

var (
	ErrNoResults           = errors.New("no results to score")
	ErrRoundResultIncomplete = errors.New("round result is incomplete")
)

type TournamentScoreId string

type PlayedResult struct {
	DealNo        int
	NsTeamId      string
	EwTeamId      string
	Contract      string
	Tricks        int
	Declarer      string
	NsVulnerable  bool
	EwVulnerable  bool
}

type BoardScore struct {
	DealNo   int
	NsTeamId string
	EwTeamId string
	NsPoints int
	NsAward  float64
	EwAward  float64
}

type Standing struct {
	TeamId     string
	TotalScore float64
	Percentage float64
	Rank       int
}

type TournamentScoreState struct {
	Id          TournamentScoreId
	Method      ScoringMethodName
	BoardScores []BoardScore
	Standings   []Standing
}

type TournamentScoreCalculated struct {
	TournamentId TournamentScoreId
	Method       ScoringMethodName
	BoardScores  []BoardScore
	Standings    []Standing
}

type TournamentScore struct {
	State  TournamentScoreState
	events []any
}

func CalculateTournamentScore(id TournamentScoreId, method ScoringMethod, results []PlayedResult) (*TournamentScore, error) {
	if len(results) == 0 {
		return nil, ErrNoResults
	}

	byDeal := map[int][]PlayedResult{}
	dealOrder := []int{}
	for _, result := range results {
		if result.Contract == "" {
			return nil, ErrRoundResultIncomplete
		}
		if _, ok := byDeal[result.DealNo]; !ok {
			dealOrder = append(dealOrder, result.DealNo)
		}
		byDeal[result.DealNo] = append(byDeal[result.DealNo], result)
	}
	sort.Ints(dealOrder)

	boardScores := []BoardScore{}
	totals := map[string]float64{}
	var maxPossible float64

	for _, dealNo := range dealOrder {
		dealResults := byDeal[dealNo]
		travellers := make([]BoardTraveller, 0, len(dealResults))
		for _, result := range dealResults {
			nsPoints, err := BridgeScore(result.Contract, result.Tricks, result.Declarer, result.NsVulnerable, result.EwVulnerable)
			if err != nil {
				return nil, err
			}
			travellers = append(travellers, BoardTraveller{
				NsTeamId: result.NsTeamId,
				EwTeamId: result.EwTeamId,
				NsPoints: nsPoints,
			})
		}

		awards := method.ScoreBoard(travellers)
		if len(awards) > 0 {
			maxPossible += awards[0].NsAward + awards[0].EwAward
		}
		for _, award := range awards {
			boardScores = append(boardScores, BoardScore{
				DealNo:   dealNo,
				NsTeamId: award.NsTeamId,
				EwTeamId: award.EwTeamId,
				NsPoints: award.NsPoints,
				NsAward:  award.NsAward,
				EwAward:  award.EwAward,
			})
			totals[award.NsTeamId] += award.NsAward
			totals[award.EwTeamId] += award.EwAward
		}
	}

	standings := make([]Standing, 0, len(totals))
	for teamId, total := range totals {
		percentage := 0.0
		if maxPossible > 0 {
			percentage = 100.0 * total / maxPossible
		}
		standings = append(standings, Standing{
			TeamId:     teamId,
			TotalScore: total,
			Percentage: percentage,
		})
	}

	sort.Slice(standings, func(i, j int) bool {
		if standings[i].TotalScore == standings[j].TotalScore {
			return standings[i].TeamId < standings[j].TeamId
		}
		return standings[i].TotalScore > standings[j].TotalScore
	})

	// Dense rank by total score.
	for i := range standings {
		if i > 0 && standings[i].TotalScore == standings[i-1].TotalScore {
			standings[i].Rank = standings[i-1].Rank
		} else {
			standings[i].Rank = i + 1
		}
	}

	score := &TournamentScore{
		State: TournamentScoreState{
			Id:          id,
			Method:      method.Name(),
			BoardScores: boardScores,
			Standings:   standings,
		},
	}
	score.events = append(score.events, TournamentScoreCalculated{
		TournamentId: id,
		Method:       method.Name(),
		BoardScores:  append([]BoardScore(nil), boardScores...),
		Standings:    append([]Standing(nil), standings...),
	})
	return score, nil
}

func (t *TournamentScore) GetEvents() []any {
	return t.events
}

func (t *TournamentScore) Commit() {
	t.events = nil
}

type TournamentScoreRepository interface {
	Save(score *TournamentScore) error
	Load(id *TournamentScoreId, method ScoringMethodName) (*TournamentScore, error)
}

type PlayedResultRepository interface {
	Replace(id TournamentScoreId, results []PlayedResult) error
	Find(id TournamentScoreId) ([]PlayedResult, error)
}

type StandingDto struct {
	TeamId       string
	TeamNumber   int
	TeamName     string
	TotalScore   float64
	Percentage   float64
	Rank         int
}

type BoardScoreDto struct {
	DealNo       int
	NsTeamId     string
	NsTeamNumber int
	EwTeamId     string
	EwTeamNumber int
	NsPoints     int
	NsAward      float64
	EwAward      float64
}

type TournamentScoreReadRepository interface {
	FindStandings(tournamentId string, method ScoringMethodName) ([]StandingDto, error)
	FindBoardScores(tournamentId string, method ScoringMethodName, dealNo *int) ([]BoardScoreDto, error)
}
