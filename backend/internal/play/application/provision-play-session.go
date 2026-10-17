package play

import (
	"errors"
	"slices"

	"bridge-tab/internal/idutil"
	domain "bridge-tab/internal/play/domain"
)

var ErrPlaySessionAlreadyProvisioned = errors.New("play session is already provisioned")
var ErrMisconfiguredPlay = errors.New("play is not configured correctly")

type TeamDto struct {
	Id      string
	Number  int
	Players []string
}

type TableScheduleBoardDto struct {
	DealNo     int
	Vulnerable string
}

type TableScheduleDto struct {
	RoundNo     int
	TableNumber int
	NsTeamId    string
	EwTeamId    string
	SetId       string
	SetLabel    string
	Boards      []TableScheduleBoardDto
}

type ProvisionPlaySessionCommand struct {
	PlaySessionId  string
	Teams          []TeamDto
	TableSchedules []TableScheduleDto
}

func (c *ProvisionPlaySessionCommand) Execute(repository domain.PlaySessionRepository) error {
	playSessionId := domain.PlaySessionId(c.PlaySessionId)

	existing, err := repository.Load(playSessionId)
	if err != nil && !errors.Is(err, domain.ErrPlaySessionNotProvisioned) {
		return err
	}
	if existing != nil {
		return ErrPlaySessionAlreadyProvisioned
	}

	teams := make([]*domain.Team, len(c.Teams))
	for i, team := range c.Teams {
		players := make([]domain.PlayerId, len(team.Players))
		for j, memberId := range team.Players {
			playerId := domain.PlayerId(memberId)
			players[j] = playerId
		}
		teamId := domain.TeamId(team.Id)
		teams[i] = &domain.Team{Id: teamId, Number: team.Number, Players: players}
	}

	tableSchedules := make([]*domain.TableSchedule, len(c.TableSchedules))
	for i, schedule := range c.TableSchedules {
		nsTeamId := domain.TeamId(schedule.NsTeamId)
		ewTeamId := domain.TeamId(schedule.EwTeamId)
		nsIndex := slices.IndexFunc(teams, func(team *domain.Team) bool {
			return idutil.SameId(team.Id, nsTeamId)
		})
		ewIndex := slices.IndexFunc(teams, func(team *domain.Team) bool {
			return idutil.SameId(team.Id, ewTeamId)
		})
		if nsIndex == -1 || ewIndex == -1 {
			return ErrMisconfiguredPlay
		}
		boards := make([]domain.TableScheduleBoard, len(schedule.Boards))
		for j, board := range schedule.Boards {
			boards[j] = domain.TableScheduleBoard{DealNo: board.DealNo, Vulnerable: board.Vulnerable}
		}
		tableSchedules[i] = &domain.TableSchedule{
			RoundNo:     schedule.RoundNo,
			TableNumber: schedule.TableNumber,
			NsTeamId:    teams[nsIndex].Id,
			EwTeamId:    teams[ewIndex].Id,
			SetId:       schedule.SetId,
			SetLabel:    schedule.SetLabel,
			Boards:      boards,
		}
	}

	session, err := domain.ProvisionPlaySession(playSessionId, teams, tableSchedules)
	if err != nil {
		return err
	}
	return repository.Save(session)
}
