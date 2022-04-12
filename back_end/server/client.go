package server

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type playerMessageType string

const (
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	joinRoom      playerMessageType = "join_room"
	setPlayerName playerMessageType = "set_player_name"
)

type PlayerMessage struct {
	MessageType playerMessageType `json:"messageType"`
	Content     string            `json:"content"`
}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	Hub          *Hub
	PlayerName   string
	Conn         *websocket.Conn
	ReceiveChan  chan []byte
	Room         *Room
	JoinRoomChan chan *Room
}

func ServerWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{
		Hub:          hub,
		PlayerName:   "",
		Conn:         conn,
		ReceiveChan:  make(chan []byte),
		JoinRoomChan: make(chan *Room),
	}
	log.Println(client.Room)
	go client.tickerHandler()
	go client.playerMessageHandler()
	go client.serverMessageHandler()
}
