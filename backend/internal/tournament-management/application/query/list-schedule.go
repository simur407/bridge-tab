package tournament_management

import (
	domain "bridge-tab/internal/tournament-management/domain"
)

type ListScheduleQuery struct {
	TournamentId string
}

func (q *ListScheduleQuery) Execute(repo domain.ScheduleReadRepository) ([]domain.ScheduleDto, error) {
	return repo.FindAll(&q.TournamentId)
}
