package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
)

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

type RoomDetail struct {
	Position    int32
	HandlerChan chan []byte
}

type Client struct {
	Hub             *Hub
	PlayerName      string
	Conn            *websocket.Conn
	Room            *RoomDetail
	Game            *Game
	ReceiveGameChan chan *Game
}

func ServerWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{
		Hub:             hub,
		PlayerName:      "",
		Conn:            conn,
		ReceiveGameChan: make(chan *Game),
	}
	go client.tickerHandler()
	go client.playerMessageHandler()
	go client.receiveGameHandler()
}

func (c *Client) sendRoomList() {
	content, _ := json.Marshal(c.Hub.roomList())
	roomListMessage := ResponseMessage{
		MessageType: roomList,
		Content:     string(content),
	}
	sendMessage, _ := json.Marshal(roomListMessage)
	c.Conn.WriteMessage(websocket.TextMessage, sendMessage)
}

func (c *Client) sendJoinRoomFail(errMsg string) {
	roomListMessage := ResponseMessage{
		MessageType: joinRoomFail,
		Content:     errMsg,
	}
	sendMessage, _ := json.Marshal(roomListMessage)
	c.Conn.WriteMessage(websocket.TextMessage, sendMessage)
}
