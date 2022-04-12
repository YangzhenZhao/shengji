package server

type Room struct {
	ID      string
	Clients []*Client
}
