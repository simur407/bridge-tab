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

	boardNo, err := strconv.Atoi(strings.TrimSpace(c.FormValue("board")))
	if err != nil || boardNo < 1 {
		return reject(c, "protocols", data, fmt.Errorf("zły numer rozdania"))
	}
	vulnerable, err := vulnerableValue(c.FormValue("vulnerable"))
	if err != nil {
		return reject(c, "protocols", data, err)
	}
	cmd := tournament_cmd.CreateBoardProtocol{
		TournamentId: tournament.Id,
		BoardNo:      boardNo,
		Vulnerable:   int(vulnerable),
	}
	if err := cmd.Execute(s.tournament); err != nil {
		return reject(c, "protocols", data, err)
	}
	return c.Redirect(path("/tournaments/"+tournament.Id+"/protocols"), fiber.StatusSeeOther)
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
	return c.Redirect(path("/tournaments/"+tournament.Id+"/protocols"), fiber.StatusSeeOther)
}

func protocolsPage(s stores, tournament *tournament_domain.TournamentDto) (fiber.Map, error) {
	protocols, err := (&tournament_query.ListBoardProtocolsQuery{TournamentId: tournament.Id}).Execute(s.boards)
	if err != nil {
		return nil, err
	}
	sort.Slice(protocols, func(i, j int) bool {
		return protocols[i].BoardNo < protocols[j].BoardNo
	})
	rows := make([]protocolRow, len(protocols))
	for i, protocol := range protocols {
		rows[i] = protocolRow{
			BoardNo:    protocol.BoardNo,
			Vulnerable: vulLabel(protocol.Vulnerable),
		}
	}
	data := pageData(tournament)
	data["Protocols"] = rows
	data["FormBoard"] = ""
	data["FormVulnerable"] = ""
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
