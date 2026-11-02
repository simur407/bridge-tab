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
	HandsSet   bool
	North      string
	East       string
	South      string
	West       string
}

func getProtocols(c *fiber.Ctx) error {
	s := newStores(c)
	tournament, err := loadTournament(c, s)
	if err != nil {
		return err
	}
	data, err := protocolsPage(s, tournament, c.Query("reveal") == "1")
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
	reveal := c.FormValue("reveal") == "1" || c.Query("reveal") == "1"
	data, err := protocolsPage(s, tournament, reveal)
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
	return c.Redirect(protocolsPath(tournament.Id, reveal), fiber.StatusSeeOther)
}

func postRemoveProtocol(c *fiber.Ctx) error {
	s := newStores(c)
	tournament, err := loadTournament(c, s)
	if err != nil {
		return err
	}
	reveal := c.FormValue("reveal") == "1" || c.Query("reveal") == "1"
	data, err := protocolsPage(s, tournament, reveal)
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
	return c.Redirect(protocolsPath(tournament.Id, reveal), fiber.StatusSeeOther)
}

func postProtocolHands(c *fiber.Ctx) error {
	s := newStores(c)
	tournament, err := loadTournament(c, s)
	if err != nil {
		return err
	}
	reveal := c.FormValue("reveal") == "1"
	data, err := protocolsPage(s, tournament, reveal)
	if err != nil {
		return err
	}
	boardNo, err := pathInt(c, "board")
	if err != nil {
		return err
	}

	north := strings.TrimSpace(c.FormValue("north"))
	east := strings.TrimSpace(c.FormValue("east"))
	south := strings.TrimSpace(c.FormValue("south"))
	west := strings.TrimSpace(c.FormValue("west"))
	data["FormHandsBoard"] = boardNo
	data["FormNorth"] = north
	data["FormEast"] = east
	data["FormSouth"] = south
	data["FormWest"] = west

	cmd := tournament_cmd.SetBoardHandsCommand{
		TournamentId: tournament.Id,
		BoardNo:      boardNo,
		North:        north,
		East:         east,
		South:        south,
		West:         west,
		Overwrite:    true,
	}
	if err := cmd.Execute(s.tournament); err != nil {
		return reject(c, "protocols", data, err)
	}
	return c.Redirect(protocolsPath(tournament.Id, reveal), fiber.StatusSeeOther)
}

func protocolsPage(s stores, tournament *tournament_domain.TournamentDto, reveal bool) (fiber.Map, error) {
	protocols, err := (&tournament_query.ListBoardProtocolsQuery{TournamentId: tournament.Id}).Execute(s.boards)
	if err != nil {
		return nil, err
	}
	sort.Slice(protocols, func(i, j int) bool {
		return protocols[i].BoardNo < protocols[j].BoardNo
	})
	showHands := tournament.FinishedAt != "" || reveal
	rows := make([]protocolRow, len(protocols))
	for i, protocol := range protocols {
		row := protocolRow{
			BoardNo:    protocol.BoardNo,
			Vulnerable: vulLabel(protocol.Vulnerable),
			HandsSet:   protocol.HandsComplete(),
		}
		if showHands {
			row.North = protocol.North
			row.East = protocol.East
			row.South = protocol.South
			row.West = protocol.West
		}
		rows[i] = row
	}
	data := pageData(tournament)
	data["Protocols"] = rows
	data["FormBoard"] = ""
	data["FormVulnerable"] = ""
	data["FormHandsBoard"] = 0
	data["FormNorth"] = ""
	data["FormEast"] = ""
	data["FormSouth"] = ""
	data["FormWest"] = ""
	data["ShowHands"] = showHands
	data["CanHideHands"] = tournament.FinishedAt == ""
	data["Reveal"] = reveal
	return data, nil
}

func protocolsPath(tournamentId string, reveal bool) string {
	target := "/tournaments/" + tournamentId + "/protocols"
	if reveal {
		target += "?reveal=1"
	}
	return path(target)
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
