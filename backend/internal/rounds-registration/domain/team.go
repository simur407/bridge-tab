package rounds_registration

type TeamId string
type PlayerId string

type Team struct {
	Id      TeamId
	Players []PlayerId
}
