package play

import (
	"errors"
	"sort"
	"strings"

	domain "bridge-tab/internal/play/domain"
)

type ListMatchesQuery struct {
	PlaySessionId string
}

func (query *ListMatchesQuery) Execute(
	playSessionRepository domain.PlaySessionRepository,
	boardResultReadRepository domain.BoardResultReadRepository,
) ([]ScheduledMatchDto, error) {
	playSessionId := domain.PlaySessionId(query.PlaySessionId)
	session, err := playSessionRepository.Load(playSessionId)
	if errors.Is(err, domain.ErrPlaySessionNotProvisioned) {
		return []ScheduledMatchDto{}, nil
	}
	if err != nil {
		return nil, err
	}

	results, err := boardResultReadRepository.FindAll(string(playSessionId))
	if err != nil {
		return nil, err
	}
	matches := session.Matches()
	recorded := map[resultKey]domain.BoardResultDto{}
	for _, result := range results {
		recorded[resultKeyOf(result.DealNo, result.NsTeamId, result.EwTeamId)] = result
	}

	matchesDto := make([]ScheduledMatchDto, 0, len(matches))
	for _, match := range matches {
		dto := scheduledMatch(session, match)
		if result, ok := recorded[resultKeyOf(match.DealNo, string(match.NsTeamId), string(match.EwTeamId))]; ok {
			applyResult(&dto, result)
		}
		matchesDto = append(matchesDto, dto)
	}
	sort.Slice(matchesDto, func(i, j int) bool {
		if matchesDto[i].DealNo != matchesDto[j].DealNo {
			return matchesDto[i].DealNo < matchesDto[j].DealNo
		}
		if matchesDto[i].NsTeamNumber != matchesDto[j].NsTeamNumber {
			return matchesDto[i].NsTeamNumber < matchesDto[j].NsTeamNumber
		}
		return matchesDto[i].EwTeamNumber < matchesDto[j].EwTeamNumber
	})
	return matchesDto, nil
}

type resultKey struct {
	deal int
	ns   string
	ew   string
}

func resultKeyOf(deal int, ns string, ew string) resultKey {
	return resultKey{deal: deal, ns: strings.ToLower(ns), ew: strings.ToLower(ew)}
}
