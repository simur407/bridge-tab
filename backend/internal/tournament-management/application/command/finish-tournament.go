package tournament_management

import (
	"errors"

	"bridge-tab/internal/idutil"
	play_domain "bridge-tab/internal/play/domain"
	score_registration "bridge-tab/internal/score-registration/application"
	score_registration_domain "bridge-tab/internal/score-registration/domain"
	domain "bridge-tab/internal/tournament-management/domain"
)

var ErrBoardResultsIncomplete = errors.New("not all board results have been recorded")

type FinishTournamentCommand struct {
	TournamentId string
}

func (c *FinishTournamentCommand) Execute(
	tournamentRepo domain.TournamentRepository,
	playSessionRepo play_domain.PlaySessionRepository,
	boardResultReadRepo play_domain.BoardResultReadRepository,
	playedResultRepo score_registration_domain.PlayedResultRepository,
	scoreRepo score_registration_domain.TournamentScoreRepository,
) error {
	id := domain.TournamentId(c.TournamentId)
	t, err := tournamentRepo.Load(&id)
	if err != nil {
		return err
	}

	// FIXME: we need to handle this check some other way, this is a hack for now
	session, err := playSessionRepo.Load(play_domain.TournamentId(c.TournamentId))
	if err != nil {
		return err
	}
	results, err := boardResultReadRepo.FindAll(c.TournamentId)
	if err != nil {
		return err
	}
	if err := ensureBoardResults(session.State.Matches, results); err != nil {
		return err
	}

	if err := t.Finish(); err != nil {
		return err
	}

	calculate := score_registration.CalculateStandingsCommand{
		TournamentId: c.TournamentId,
		Method:       score_registration_domain.MethodMatchpoints,
	}
	if err := calculate.Execute(playedResultRepo, scoreRepo); err != nil {
		return err
	}
	return tournamentRepo.Save(t)
}

func ensureBoardResults(matches []*play_domain.Match, results []play_domain.BoardResultDto) error {
	for _, match := range matches {
		if !boardResultCovers(results, match.DealNo, string(match.NsTeamId), string(match.EwTeamId)) {
			return ErrBoardResultsIncomplete
		}
	}
	return nil
}

func boardResultCovers(results []play_domain.BoardResultDto, dealNo int, nsTeamId string, ewTeamId string) bool {
	for _, result := range results {
		if result.DealNo == dealNo && idutil.SameId(result.NsTeamId, nsTeamId) && idutil.SameId(result.EwTeamId, ewTeamId) {
			return true
		}
	}
	return false
}
