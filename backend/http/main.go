package main

import (
	"bridge-tab/http/middleware"
	auth "bridge-tab/internal/auth"
	rounds_registration_cmd "bridge-tab/internal/rounds-registration/application"
	rounds_registration_infra "bridge-tab/internal/rounds-registration/infrastructure"
	tournament_management_cmd "bridge-tab/internal/tournament-management/application/command"
	tournament_management_query "bridge-tab/internal/tournament-management/application/query"
	tournament_management_domain "bridge-tab/internal/tournament-management/domain"
	tournament_management_infra "bridge-tab/internal/tournament-management/infrastructure"
	users "bridge-tab/internal/user/application"
	users_infra "bridge-tab/internal/user/infrastructure"
	"database/sql"
	"os"
	"slices"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func main() {
	app := fiber.New(fiber.Config{
		Views: html.New("http/frontend", ".html"),
	})
	app.Use(logger.New())
	app.Use(recover.New())

	dbString := os.Getenv("DATABASE_STRING")
	db, err := sql.Open("postgres", dbString)

	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic("failed to connect to database")
	}

	tournament_management_infra.Migrate(db)
	users_infra.Migrate(db)
	rounds_registration_infra.Migrate(db)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", nil)
	})
	app.Get("/register", GetRegister)
	app.Post("/register", middleware.Transaction(db, nil), PostRegister)
	app.Get("/tournaments/:tournamentId",
		middleware.JwtGuard(),
		middleware.Transaction(db, &sql.TxOptions{ReadOnly: true}),
		GetTournament,
	)
	app.Post("/tournaments/:tournamentId/join",
		middleware.JwtGuard(),
		middleware.Transaction(db, nil),
		JoinTournament,
	)
	app.Get("/tournaments/:tournamentId/teams",
		middleware.JwtGuard(),
		middleware.Transaction(db, nil),
		GetTournamentTeams,
	)
	app.Post("/tournaments/:tournamentId/teams/:teamId/join",
		middleware.JwtGuard(),
		middleware.Transaction(db, nil),
		JoinTeam,
	)
	app.Post("/tournaments/:tournamentId/teams/:teamId/leave",
		middleware.JwtGuard(),
		middleware.Transaction(db, nil),
		LeaveTeam,
	)
	app.Get("/round-registration/:gameSessionId/add-round",
		middleware.JwtGuard(),
		middleware.Transaction(db, nil),
		GetAddRoundForm,
	)
	app.Post("/round-registration/:gameSessionId/add-round",
		middleware.JwtGuard(),
		middleware.Transaction(db, nil),
		SubmitRound,
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	app.Listen(":" + port)
}

func GetRegister(c *fiber.Ctx) error {
	return c.Render("register", fiber.Map{
		"Title":    "Rejestracja",
		"Redirect": c.Query("redirect"),
	}, "layout")
}

type RegisterUserRequestDto struct {
	Name string `json:"name"`
}

func PostRegister(c *fiber.Ctx) error {
	body := new(RegisterUserRequestDto)

	if err := c.BodyParser(body); err != nil {
		log.Debug(err)
		return err
	}

	tx := middleware.GetTransaction(c)
	repository := users_infra.PostgresUserRepository{Ctx: c.UserContext(), Tx: tx}

	id := uuid.New().String()

	command := &users.RegisterUserCommand{
		Id:   id,
		Name: body.Name,
	}
	if err := command.Execute(&repository); err != nil {
		log.Debug(err)
		return err
	}

	token, err := auth.Generate(id)

	if err != nil {
		log.Debug(err)
		return err
	}

	cookie := new(fiber.Cookie)
	cookie.Name = "token"
	cookie.Value = token
	cookie.Expires = time.Now().Add(auth.EXPIRES_AT)
	c.Cookie(cookie)

	// TODO: handle no redirect
	redirect := c.Query("redirect")
	return c.Redirect(redirect)
}

func GetTournament(c *fiber.Ctx) error {
	tournamentId := c.Params("tournamentId")

	if tournamentId == "" {
		log.Debug("tournamentId is empty")
		return c.Render("404", nil)
	}

	if err := uuid.Validate(tournamentId); err != nil {
		log.Debug(err)
		return c.Render("404", nil)
	}

	joinedTournament := c.Cookies("tournamentId")

	if joinedTournament != "" {
		return c.Redirect("/tournaments/" + joinedTournament + "/teams")
	}

	getTournament := tournament_management_query.GetTournamentById{
		Id: tournamentId,
	}
	tournament, err := getTournament.Execute(&tournament_management_infra.PostgresTournamentReadRepository{
		Ctx: c.UserContext(),
		Tx:  middleware.GetTransaction(c),
	})

	if err != nil {
		log.Debug(err)
		return notFound(c, "Błąd przy pobieraniu turnieju")
	}

	if tournament == nil {
		log.Debug("tournament not found")
		return notFound(c, "Nie znaleziono turnieju")
	}

	if tournament.StartedAt != "" {
		log.Debug("tournament is already started")
		return notFound(c, "Nie można już dołączać do turnieju")
	}

	return c.Render("tournament", fiber.Map{
		"Title":          tournament.Name,
		"TournamentName": tournament.Name,
		"TournamentId":   tournament.Id,
	}, "layout")
}

func JoinTournament(c *fiber.Ctx) error {
	tournamentId := c.Params("tournamentId")

	if tournamentId == "" {
		log.Debug("tournamentId is empty")
		return c.Render("404", nil)
	}

	if err := uuid.Validate(tournamentId); err != nil {
		log.Debug(err)
		return c.Render("404", nil)
	}

	contestantId := c.Locals("user").(middleware.UserMetadata).Id

	joinTournament := tournament_management_cmd.JoinTournamentCommand{
		TournamentId: tournamentId,
		ContestantId: contestantId,
	}
	err := joinTournament.Execute(&tournament_management_infra.PostgresTournamentRepository{
		Ctx: c.UserContext(),
		Tx:  middleware.GetTransaction(c),
	})

	if err != nil {
		log.Debug(err)
		return c.Render("404", nil)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "tournamentId",
		Value:    tournamentId,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	})

	return c.Redirect("/tournaments/" + tournamentId + "/teams")
}

func GetTournamentTeams(c *fiber.Ctx) error {
	tournamentId := c.Params("tournamentId")

	if tournamentId == "" {
		log.Debug("tournamentId is empty")
		return notFound(c, "Nie znaleziono turnieju")
	}

	if err := uuid.Validate(tournamentId); err != nil {
		log.Debug(err)
		return notFound(c, "Niepoprawny identyfikator turnieju")
	}

	getTournament := tournament_management_query.GetTournamentById{
		Id: tournamentId,
	}
	tournament, err := getTournament.Execute(&tournament_management_infra.PostgresTournamentReadRepository{
		Ctx: c.UserContext(),
		Tx:  middleware.GetTransaction(c),
	})

	if err != nil {
		log.Debug(err)
		return notFound(c, "Błąd przy pobieraniu turnieju")
	}

	if tournament == nil {
		log.Debug("tournament not found")
		return notFound(c, "Nie znaleziono turnieju")
	}

	joinedTeamId := c.Cookies("teamId", "")

	getTournamentTeams := tournament_management_query.ListTeamsQuery{
		TournamentId: tournamentId,
	}
	teams, err := getTournamentTeams.Execute(&tournament_management_infra.PostgresTeamReadRepository{
		Ctx: c.UserContext(),
		Tx:  middleware.GetTransaction(c),
	})

	joinedTeamIndex := slices.IndexFunc(teams, func(team tournament_management_domain.TeamDto) bool {
		return team.Id == joinedTeamId
	})

	var joinedTeam tournament_management_domain.TeamDto
	if joinedTeamIndex != -1 {
		joinedTeam = teams[joinedTeamIndex]
	}

	if err != nil {
		log.Debug(err)
		return notFound(c, "Błąd przy pobieraniu drużyn")
	}

	return c.Render("teams", fiber.Map{
		"Title":          tournament.Name,
		"TournamentId":   tournament.Id,
		"TournamentName": tournament.Name,
		"Teams":          teams,
		"JoinedTeam":     joinedTeam,
		"Started":        tournament.StartedAt != "",
		"UserId":         c.Locals("user").(middleware.UserMetadata).Id,
	}, "layout")
}

func JoinTeam(c *fiber.Ctx) error {
	tournamentId := c.Params("tournamentId")

	if tournamentId == "" {
		log.Debug("tournamentId is empty")
		return c.Render("404", nil)
	}
	if err := uuid.Validate(tournamentId); err != nil {
		log.Debug(err)
		return c.Render("404", nil)
	}

	teamId := c.Params("teamId")

	if teamId == "" {
		log.Debug("teamId is empty")
		return c.Render("404", nil)
	}
	if err := uuid.Validate(teamId); err != nil {
		log.Debug(err)
		return c.Render("404", nil)
	}

	contestantId := c.Locals("user").(middleware.UserMetadata).Id

	joinTeam := tournament_management_cmd.JoinTeamCommand{
		TournamentId: tournamentId,
		TeamId:       teamId,
		ContestantId: contestantId,
	}
	err := joinTeam.Execute(&tournament_management_infra.PostgresTournamentRepository{
		Ctx: c.UserContext(),
		Tx:  middleware.GetTransaction(c),
	})

	if err != nil {
		log.Debug(err)
		return c.Render("404", nil)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "teamId",
		Value:    teamId,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	})

	return c.Redirect("/tournaments/" + tournamentId + "/teams")
}

func LeaveTeam(c *fiber.Ctx) error {
	tournamentId := c.Params("tournamentId")

	if tournamentId == "" {
		log.Debug("tournamentId is empty")
		return c.Render("404", nil)
	}
	if err := uuid.Validate(tournamentId); err != nil {
		log.Debug(err)
		return c.Render("404", nil)
	}

	teamId := c.Params("teamId")

	if teamId == "" {
		log.Debug("teamId is empty")
		return c.Render("404", nil)
	}
	if err := uuid.Validate(teamId); err != nil {
		log.Debug(err)
		return c.Render("404", nil)
	}

	contestantId := c.Locals("user").(middleware.UserMetadata).Id

	leaveTeam := tournament_management_cmd.LeaveTeamCommand{
		TournamentId: tournamentId,
		TeamId:       teamId,
		ContestantId: contestantId,
	}
	err := leaveTeam.Execute(&tournament_management_infra.PostgresTournamentRepository{
		Ctx: c.UserContext(),
		Tx:  middleware.GetTransaction(c),
	})

	if err != nil {
		log.Debug(err)
		return c.Render("404", nil)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "teamId",
		Expires:  time.Now(),
		HTTPOnly: true,
	})

	return c.Redirect("/tournaments/" + tournamentId + "/teams")
}

func GetAddRoundForm(c *fiber.Ctx) error {
	gameSessionId := c.Params("gameSessionId")
	if gameSessionId == "" {
		log.Debug("gameSessionId is empty")
		return notFound(c, "Nie znaleziono rozgrywki")
	}
	if err := uuid.Validate(gameSessionId); err != nil {
		log.Debug(err)
		return notFound(c, "Niepoprawny identyfikator rozgrywki")
	}

	playerId := c.Locals("user").(middleware.UserMetadata).Id
	getPlayerTeam := tournament_management_query.GetTeamByMemberQuery{TournamentId: gameSessionId, MemberId: playerId}
	playerTeam, err := getPlayerTeam.Execute(&tournament_management_infra.PostgresTeamReadRepository{
		Ctx: c.UserContext(),
		Tx:  middleware.GetTransaction(c),
	})

	if err != nil {
		log.Debug(err)
		return notFound(c, "Nie dołączono do żadnego zespołu")
	}

	var success string
	if c.Query("success") != "" {
		success = "Pomyślnie dodano rundę"
	}

	return c.Render("add-round", fiber.Map{
		"Success":       success,
		"Title":         "Dodaj rundę",
		"GameSessionId": c.Params("gameSessionId"),
		"PlayerTeam":    playerTeam.Name,
		"UserId":        playerId,
	}, "layout")
}

type Round struct {
	DealNo         int    `json:"dealNo"`
	VersusTeamName string `json:"versusTeamName"`
	Contract       string `json:"contract"`
	Tricks         int    `json:"tricks"`
	Declarer       string `json:"declarer"`
	OpeningLead    string `json:"openingLead"`
}

func SubmitRound(c *fiber.Ctx) error {
	gameSessionId := c.Params("gameSessionId")

	if gameSessionId == "" {
		log.Debug("gameSessionId is empty")
		return notFound(c, "Nie znaleziono rozgrywki")
	}
	if err := uuid.Validate(gameSessionId); err != nil {
		log.Debug(err)
		return notFound(c, "Niepoprawny identyfikator rozgrywki")
	}

	body := new(Round)
	if err := c.BodyParser(body); err != nil {
		log.Debug(err)
		return notFound(c, "Błąd przy przetwarzaniu formularza")
	}

	contestantId := c.Locals("user").(middleware.UserMetadata).Id

	submitRound := rounds_registration_cmd.PlayRoundCommand{
		GameSessionId:  gameSessionId,
		PlayerId:       contestantId,
		DealNo:         body.DealNo,
		VersusTeamName: body.VersusTeamName,
		Contract:       body.Contract,
		Tricks:         body.Tricks,
		Declarer:       body.Declarer,
		OpeningLead:    body.OpeningLead,
	}
	err := submitRound.Execute(&rounds_registration_infra.PostgresGameSessionRepository{
		Ctx: c.UserContext(),
		Tx:  middleware.GetTransaction(c),
	}, &tournament_management_infra.PostgresTeamReadRepository{
		Ctx: c.UserContext(),
		Tx:  middleware.GetTransaction(c),
	})

	if err != nil {
		log.Debug(err)
		return c.Render("error", fiber.Map{
			"Error": err.Error(),
		})
	}

	c.Append("HX-Redirect", "/round-registration/"+gameSessionId+"/add-round?success=true")
	return c.SendStatus(302)
}

func notFound(c *fiber.Ctx, message string) error {
	return c.Render("404", fiber.Map{
		"Title":   "Błąd",
		"Message": message,
	}, "layout")
}
