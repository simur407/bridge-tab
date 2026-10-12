package play

import (
	"errors"
	"sort"
	"strings"

	domain "bridge-tab/internal/play/domain"
)

type ListMatchesQuery struct {
	TournamentId string
}

func (query *ListMatchesQuery) Execute(
	playSessionRepository domain.PlaySessionRepository,
	boardResultReadRepository domain.BoardResultReadRepository,
) ([]ScheduledMatchDto, error) {
	session, err := playSessionRepository.Load(domain.TournamentId(query.TournamentId))
	if errors.Is(err, domain.ErrPlaySessionNotProvisioned) {
		return []ScheduledMatchDto{}, nil
	}
	if err != nil {
		return nil, err
	}

	results, err := boardResultReadRepository.FindAll(query.TournamentId)
	if err != nil {
		return nil, err
	}
	recorded := map[resultKey]domain.BoardResultDto{}
	for _, result := range results {
		recorded[resultKeyOf(result.DealNo, result.NsTeamId, result.EwTeamId)] = result
	}

	matches := make([]ScheduledMatchDto, 0, len(session.State.Matches))
	for _, match := range session.State.Matches {
		dto := scheduledMatch(session, match)
		if result, ok := recorded[resultKeyOf(match.DealNo, string(match.NsTeamId), string(match.EwTeamId))]; ok {
			applyResult(&dto, result)
		}
		matches = append(matches, dto)
	}
	sort.Slice(matches, func(i, j int) bool {
		if matches[i].DealNo != matches[j].DealNo {
			return matches[i].DealNo < matches[j].DealNo
		}
		if matches[i].NsTeamNumber != matches[j].NsTeamNumber {
			return matches[i].NsTeamNumber < matches[j].NsTeamNumber
		}
		return matches[i].EwTeamNumber < matches[j].EwTeamNumber
	})
	return matches, nil
}

type resultKey struct {
	deal int
	ns   string
	ew   string
}

func resultKeyOf(deal int, ns string, ew string) resultKey {
	return resultKey{deal: deal, ns: strings.ToLower(ns), ew: strings.ToLower(ew)}
}
