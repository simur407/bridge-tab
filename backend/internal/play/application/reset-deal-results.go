package play

import (
	"errors"

	domain "bridge-tab/internal/play/domain"
	score_app "bridge-tab/internal/score-registration/application"
	score_domain "bridge-tab/internal/score-registration/domain"
)

type ResetDealResultsCommand struct {
	PlaySessionId string
	DealNo        int
}

func (c *ResetDealResultsCommand) Execute(
	boardResultRepository domain.BoardResultRepository,
	boardResultReadRepository domain.BoardResultReadRepository,
	playedResultRepository score_domain.PlayedResultRepository,
	scoreRepository score_domain.TournamentScoreRepository,
) error {
	if c.PlaySessionId == "" {
		return errors.New("play session id is empty")
	}
	if c.DealNo < 1 {
		return errors.New("deal number is empty")
	}

	results, err := boardResultReadRepository.FindAll(c.PlaySessionId)
	if err != nil {
		return err
	}

	playSessionId := domain.PlaySessionId(c.PlaySessionId)
	for _, dto := range results {
		if dto.DealNo != c.DealNo {
			continue
		}
		result, err := boardResultRepository.Load(
			playSessionId,
			dto.DealNo,
			domain.TeamId(dto.NsTeamId),
			domain.TeamId(dto.EwTeamId),
		)
		if err != nil {
			return err
		}
		result.Reset()
		if err := boardResultRepository.Save(result); err != nil {
			return err
		}
		if err := unregisterPlayedResult(c.PlaySessionId, result.State, playedResultRepository); err != nil {
			return err
		}
	}

	calc := score_app.CalculateStandingsCommand{
		TournamentId: c.PlaySessionId,
		Method:       score_domain.MethodMatchpoints,
	}
	err = calc.Execute(playedResultRepository, scoreRepository)
	if errors.Is(err, score_app.ErrPlayedResultsNotFound) {
		return scoreRepository.Delete(
			score_domain.TournamentScoreId(c.PlaySessionId),
			score_domain.MethodMatchpoints,
		)
	}
	return err
}
