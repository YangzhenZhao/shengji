package server

import (
	"bytes"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func (c *Client) playerMessageHandler() {
	defer func() {
		c.Conn.Close()
	}()
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		playerMessage := RequestMessage{}
		err = json.Unmarshal(message, &playerMessage)
		if err != nil {
			log.Printf("unmarshal playerMessage err: %+v", err)
			continue
		}
		switch playerMessage.MessageType {
		case joinRoom:
			c.Hub.JoinRoomRequestChan <- &JoinRoomRequest{
				PlayerName: c.PlayerName,
				RoomID:     playerMessage.Content,
				Conn:       c.Conn,
			}
		case setPlayerName:
			c.PlayerName = string(playerMessage.Content)
			c.Hub.RegisterClientChan <- &RegisterClientRequest{PlayerName: c.PlayerName, Client: c}
			c.sendRoomList()
		default:
			log.Println("用户信息格式错误!")
		}
	}
}

func (c *Client) tickerHandler() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()
	for {
		<-ticker.C
		c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
		if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			return
		}
	}
}
