package play

import (
	"errors"
	"slices"
)

type PlaySessionId string
type TeamId string

type BoardResultState struct {
	PlaySessionId PlaySessionId
	DealNo        int
	NsTeamId      TeamId
	EwTeamId      TeamId
	Contract      string
	Tricks        int
	Declarer      string
	OpeningLead   string
}

type BoardResultRecorded struct {
	PlaySessionId PlaySessionId
	DealNo        int
	NsTeamId      TeamId
	EwTeamId      TeamId
	Contract      string
	Tricks        int
	Declarer      string
	OpeningLead   string
}

type BoardResultEdited struct {
	PlaySessionId PlaySessionId
	DealNo        int
	NsTeamId      TeamId
	EwTeamId      TeamId
	Contract      string
	Tricks        int
	Declarer      string
	OpeningLead   string
}

var ErrBoardResultNotFound = errors.New("board result not found")
var ErrBoardResultAlreadyRecorded = errors.New("board result already recorded")

type BoardResult struct {
	State  BoardResultState
	events []any
}

func RecordBoardResult(playSessionId PlaySessionId, dealNo int, nsTeamId TeamId, ewTeamId TeamId, contract string, tricks int, declarer string, openingLead string) *BoardResult {
	state := BoardResultState{
		PlaySessionId: playSessionId,
		DealNo:        dealNo,
		NsTeamId:      nsTeamId,
		EwTeamId:      ewTeamId,
		Contract:      contract,
		Tricks:        tricks,
		Declarer:      declarer,
		OpeningLead:   openingLead,
	}
	return &BoardResult{
		State: state,
		events: []any{BoardResultRecorded{
			PlaySessionId: playSessionId,
			DealNo:        dealNo,
			NsTeamId:      nsTeamId,
			EwTeamId:      ewTeamId,
			Contract:      contract,
			Tricks:        tricks,
			Declarer:      declarer,
			OpeningLead:   openingLead,
		}},
	}
}

func (b *BoardResult) Edit(contract string, tricks int, declarer string, openingLead string) {
	b.State.Contract = contract
	b.State.Tricks = tricks
	b.State.Declarer = declarer
	b.State.OpeningLead = openingLead

	b.events = append(b.events, BoardResultEdited{
		PlaySessionId: b.State.PlaySessionId,
		DealNo:        b.State.DealNo,
		NsTeamId:      b.State.NsTeamId,
		EwTeamId:      b.State.EwTeamId,
		Contract:      contract,
		Tricks:        tricks,
		Declarer:      declarer,
		OpeningLead:   openingLead,
	})
}

func (b *BoardResult) GetEvents() []any {
	return b.events
}

func (b *BoardResult) Commit() {
	b.events = slices.Delete(b.events, 0, len(b.events))
}

type BoardResultRepository interface {
	Save(result *BoardResult) error
	Load(playSessionId PlaySessionId, dealNo int, teamA TeamId, teamB TeamId) (*BoardResult, error)
}

type BoardResultDto struct {
	DealNo       int
	NsTeamId     string
	EwTeamId     string
	NsTeamNumber int
	EwTeamNumber int
	Contract     string
	Tricks       int
	Declarer     string
	OpeningLead  string
}

type BoardResultReadRepository interface {
	FindAll(playSessionId string) ([]BoardResultDto, error)
	FindByTeams(playSessionId string, dealNo int, teamAId string, teamBId string) (*BoardResultDto, error)
}
