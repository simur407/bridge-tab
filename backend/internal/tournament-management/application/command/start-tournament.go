package tournament_management

import (
	"slices"

	"bridge-tab/internal/idutil"
	play "bridge-tab/internal/play/application"
	play_domain "bridge-tab/internal/play/domain"
	domain "bridge-tab/internal/tournament-management/domain"
)

type StartTurnamentCommand struct {
	TournamentId string
}

func (c *StartTurnamentCommand) Execute(
	tournamentRepo domain.TournamentRepository,
	playSessionRepo play_domain.PlaySessionRepository,
) error {
	tournamentId := domain.TournamentId(c.TournamentId)
	t, err := tournamentRepo.Load(&tournamentId)
	if err != nil {
		return err
	}

	if err := t.Start(); err != nil {
		return err
	}

	if err := tournamentRepo.Save(t); err != nil {
		return err
	}

	playTeams := make([]play.TeamDto, len(t.State.Teams))
	for i, team := range t.State.Teams {
		players := make([]string, len(team.State.Members))
		for j, member := range team.State.Members {
			playerId := string(member.Id)
			players[j] = playerId
		}
		teamId := string(team.State.Id)
		playTeams[i] = play.TeamDto{Id: teamId, Number: team.State.Number, Players: players}
	}

	tableSchedules := make([]play.TableScheduleDto, 0, len(t.State.Schedule))
	for _, entry := range t.State.Schedule {
		setIndex := slices.IndexFunc(t.State.Sets, func(set *domain.Set) bool {
			return idutil.SameId(set.Id, entry.SetId)
		})
		if setIndex == -1 {
			return play.ErrMisconfiguredPlay
		}
		set := t.State.Sets[setIndex]
		tableIndex := slices.IndexFunc(t.State.Tables, func(table *domain.Table) bool {
			return idutil.SameId(table.State.Id, entry.TableId)
		})
		if tableIndex == -1 {
			return play.ErrMisconfiguredPlay
		}

		boards := make([]play.TableScheduleBoardDto, 0, len(set.BoardNos))
		for _, boardNo := range set.BoardNos {
			protocolIndex := slices.IndexFunc(t.State.BoardProtocols, func(protocol *domain.BoardProtocol) bool {
				return protocol.BoardNo == boardNo
			})
			if protocolIndex == -1 {
				return play.ErrMisconfiguredPlay
			}
			boards = append(boards, play.TableScheduleBoardDto{
				DealNo:     boardNo,
				Vulnerable: t.State.BoardProtocols[protocolIndex].Vulnerable.String(),
			})
		}

		nsTeamId := string(entry.NsTeamId)
		ewTeamId := string(entry.EwTeamId)
		setId := string(entry.SetId)
		tableSchedules = append(tableSchedules, play.TableScheduleDto{
			RoundNo:     entry.RoundNo,
			TableNumber: t.State.Tables[tableIndex].State.Number,
			NsTeamId:    nsTeamId,
			EwTeamId:    ewTeamId,
			SetId:       setId,
			SetLabel:    set.Label,
			Boards:      boards,
		})
	}

	playSessionId := c.TournamentId
	provisionPlaySession := play.ProvisionPlaySessionCommand{
		PlaySessionId:  playSessionId,
		Teams:          playTeams,
		TableSchedules: tableSchedules,
	}
	return provisionPlaySession.Execute(playSessionRepo)
}
