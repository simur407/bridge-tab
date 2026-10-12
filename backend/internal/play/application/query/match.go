package play

import (
	"bridge-tab/internal/idutil"
	domain "bridge-tab/internal/play/domain"
)

type ScheduledMatchDto struct {
	DealNo       int
	NsTeamNumber int
	EwTeamNumber int
	Contract     string
	Tricks       int
	Declarer     string
	OpeningLead  string
}

func scheduledMatch(session *domain.PlaySession, match *domain.Match) ScheduledMatchDto {
	return ScheduledMatchDto{
		DealNo:       match.DealNo,
		NsTeamNumber: teamNumber(session, match.NsTeamId),
		EwTeamNumber: teamNumber(session, match.EwTeamId),
	}
}

func teamNumber(session *domain.PlaySession, id domain.TeamId) int {
	for _, team := range session.State.Teams {
		if idutil.SameId(team.Id, id) {
			return team.Number
		}
	}
	return 0
}

func applyResult(match *ScheduledMatchDto, result domain.BoardResultDto) {
	match.Contract = result.Contract
	match.Tricks = result.Tricks
	match.Declarer = result.Declarer
	match.OpeningLead = result.OpeningLead
}
