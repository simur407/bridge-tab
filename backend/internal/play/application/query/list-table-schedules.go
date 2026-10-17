package play

import (
	"errors"

	domain "bridge-tab/internal/play/domain"
)

type TableScheduleView struct {
	RoundNo      int
	TableNumber  int
	NsTeamNumber int
	EwTeamNumber int
	SetId        string
	SetLabel     string
	BoardNos     []int
}

type ListTableSchedulesQuery struct {
	PlaySessionId string
}

func (query *ListTableSchedulesQuery) Execute(playSessionRepository domain.PlaySessionRepository) ([]TableScheduleView, error) {
	playSessionId := domain.PlaySessionId(query.PlaySessionId)
	session, err := playSessionRepository.Load(playSessionId)
	if errors.Is(err, domain.ErrPlaySessionNotProvisioned) {
		return []TableScheduleView{}, nil
	}
	if err != nil {
		return nil, err
	}

	views := make([]TableScheduleView, 0, len(session.State.TableSchedules))
	for _, schedule := range session.State.TableSchedules {
		boardNos := make([]int, len(schedule.Boards))
		for i, board := range schedule.Boards {
			boardNos[i] = board.DealNo
		}
		views = append(views, TableScheduleView{
			RoundNo:      schedule.RoundNo,
			TableNumber:  schedule.TableNumber,
			NsTeamNumber: teamNumber(session, schedule.NsTeamId),
			EwTeamNumber: teamNumber(session, schedule.EwTeamId),
			SetId:        schedule.SetId,
			SetLabel:     schedule.SetLabel,
			BoardNos:     boardNos,
		})
	}
	return views, nil
}
