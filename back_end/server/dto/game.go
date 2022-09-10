package dto

type TeamType int

const (
	UnknownTeam TeamType = iota
	FirstTeam
	SecondTeam
)

type GameResult struct {
	Score            int
	FinalCardWinTeam TeamType
	FinalWinTeam     TeamType
	UpLevel          int
}
