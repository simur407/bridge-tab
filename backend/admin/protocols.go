package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	tournament_cmd "bridge-tab/internal/tournament-management/application/command"
	tournament_query "bridge-tab/internal/tournament-management/application/query"
	tournament_domain "bridge-tab/internal/tournament-management/domain"

	"github.com/gofiber/fiber/v2"
)

type protocolRow struct {
	BoardNo    int
	Vulnerable string
	Pairs      []string
}

func getProtocols(c *fiber.Ctx) error {
	s := newStores(c)
	tournament, err := loadTournament(c, s)
	if err != nil {
		return err
	}
	data, err := protocolsPage(s, tournament)
	if err != nil {
		return err
	}
	return view(c, "protocols", data)
}

func postProtocol(c *fiber.Ctx) error {
	s := newStores(c)
	tournament, err := loadTournament(c, s)
	if err != nil {
		return err
	}
	data, err := protocolsPage(s, tournament)
	if err != nil {
		return err
	}
	data["FormBoard"] = strings.TrimSpace(c.FormValue("board"))
	data["FormVulnerable"] = c.FormValue("vulnerable")
	data["FormPairs"] = c.FormValue("pairs")

	boardNo, err := strconv.Atoi(strings.TrimSpace(c.FormValue("board")))
	if err != nil || boardNo < 1 {
		return reject(c, "protocols", data, fmt.Errorf("zły numer rozdania"))
	}
	vulnerable, err := vulnerableValue(c.FormValue("vulnerable"))
	if err != nil {
		return reject(c, "protocols", data, err)
	}
	lines, err := parsePairLines(c.FormValue("pairs"))
	if err != nil {
		return reject(c, "protocols", data, err)
	}
	teamPairs, err := resolvePairs(s, tournament.Id, lines)
	if err != nil {
		return reject(c, "protocols", data, err)
	}
	cmd := tournament_cmd.CreateBoardProtocol{
		TournamentId: tournament.Id,
		BoardNo:      boardNo,
		Vulnerable:   int(vulnerable),
		TeamPairs:    teamPairs,
	}
	if err := cmd.Execute(s.tournament); err != nil {
		return reject(c, "protocols", data, err)
	}
	return c.Redirect("/tournaments/"+tournament.Id+"/protocols", fiber.StatusSeeOther)
}

func postRemoveProtocol(c *fiber.Ctx) error {
	s := newStores(c)
	tournament, err := loadTournament(c, s)
	if err != nil {
		return err
	}
	data, err := protocolsPage(s, tournament)
	if err != nil {
		return err
	}
	boardNo, err := pathInt(c, "board")
	if err != nil {
		return err
	}
	cmd := tournament_cmd.RemoveBoardProtocolCommand{TournamentId: tournament.Id, BoardNo: boardNo}
	if err := cmd.Execute(s.tournament); err != nil {
		return reject(c, "protocols", data, err)
	}
	return c.Redirect("/tournaments/"+tournament.Id+"/protocols", fiber.StatusSeeOther)
}

func protocolsPage(s stores, tournament *tournament_domain.TournamentDto) (fiber.Map, error) {
	protocols, err := (&tournament_query.ListBoardProtocolsQuery{TournamentId: tournament.Id}).Execute(s.boards)
	if err != nil {
		return nil, err
	}
	teams, err := (&tournament_query.ListTeamsQuery{TournamentId: tournament.Id}).Execute(s.teams)
	if err != nil {
		return nil, err
	}
	sort.Slice(protocols, func(i, j int) bool {
		return protocols[i].BoardNo < protocols[j].BoardNo
	})
	rows := make([]protocolRow, len(protocols))
	for i, protocol := range protocols {
		pairs := make([]string, len(protocol.TeamPairs))
		for j, pair := range protocol.TeamPairs {
			ns := teamCaption(teams, pair.NS, 0)
			ew := teamCaption(teams, pair.EW, 0)
			if pair.TableNumber != nil {
				pairs[j] = fmt.Sprintf("stół %d: %s – %s", *pair.TableNumber, ns, ew)
			} else {
				pairs[j] = ns + " – " + ew
			}
		}
		rows[i] = protocolRow{
			BoardNo:    protocol.BoardNo,
			Vulnerable: vulLabel(protocol.Vulnerable),
			Pairs:      pairs,
		}
	}
	data := pageData(tournament)
	data["Protocols"] = rows
	data["FormBoard"] = ""
	data["FormVulnerable"] = ""
	data["FormPairs"] = ""
	return data, nil
}

func vulnerableValue(value string) (tournament_domain.Vulnerable, error) {
	switch value {
	case "None", "":
		return tournament_domain.None, nil
	case "NS":
		return tournament_domain.NS, nil
	case "EW":
		return tournament_domain.EW, nil
	case "Both":
		return tournament_domain.Both, nil
	default:
		return 0, fmt.Errorf("zła partia")
	}
}

func resolvePairs(s stores, tournamentId string, lines []pairLine) ([]struct {
	Table *string
	NS    string
	EW    string
}, error) {
	pairs := make([]struct {
		Table *string
		NS    string
		EW    string
	}, len(lines))
	for i, line := range lines {
		ns, err := (&tournament_query.GetTeamByNumberQuery{TournamentId: tournamentId, Number: line.NS}).Execute(s.teams)
		if err != nil {
			return nil, fmt.Errorf("drużyna NS %d: %w", line.NS, explain(err))
		}
		ew, err := (&tournament_query.GetTeamByNumberQuery{TournamentId: tournamentId, Number: line.EW}).Execute(s.teams)
		if err != nil {
			return nil, fmt.Errorf("drużyna EW %d: %w", line.EW, explain(err))
		}
		pairs[i].NS = ns.Id
		pairs[i].EW = ew.Id
		if line.Table != nil {
			table, err := (&tournament_query.GetTableByNumberQuery{TournamentId: tournamentId, Number: *line.Table}).Execute(s.tables)
			if err != nil {
				return nil, fmt.Errorf("stół %d: %w", *line.Table, explain(err))
			}
			id := table.Id
			pairs[i].Table = &id
		}
	}
	return pairs, nil
}
