package server

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
)

type serverMessageType string

const (
	roomList serverMessageType = "room_list"
)

type ServerMessage struct {
	MessageType serverMessageType `json:"messageType"`
	Content     string            `json:"content"`
}

type RegisterClientRequest struct {
	PlayerName string
	Client     *Client
}

type Hub struct {
	Clients map[string]*Client
	Rooms   map[string]*Room

	ReceiveChan        chan []byte
	RegisterClientChan chan *RegisterClientRequest
}

func NewHub() *Hub {
	roomID := uuid.NewString()
	return &Hub{
		Clients: make(map[string]*Client),
		Rooms: map[string]*Room{
			roomID: {
				ID:      roomID,
				Clients: []*Client{},
			},
		},
		ReceiveChan:        make(chan []byte),
		RegisterClientChan: make(chan *RegisterClientRequest),
	}
}

func (h *Hub) roomList() []string {
	var res []string
	for roomID := range h.Rooms {
		res = append(res, roomID)
	}
	return res
}

func (h *Hub) Run() {
	for {
		select {
		case registerClientRequest := <-h.RegisterClientChan:
			h.Clients[registerClientRequest.PlayerName] = registerClientRequest.Client
			log.Printf("%s 注册成功!\n", registerClientRequest.PlayerName)
			content, _ := json.Marshal(h.roomList())
			roomListMessage := ServerMessage{
				MessageType: roomList,
				Content:     string(content),
			}
			sendMessage, _ := json.Marshal(roomListMessage)
			registerClientRequest.Client.ReceiveChan <- sendMessage
		}
	}
}
