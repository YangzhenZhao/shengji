package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func (c *Client) playerMessageHandler() {
	defer func() {
		log.Printf("%s ws 关闭\n", c.PlayerName)
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
			log.Printf("client handler 错误 err:%+v\n", err)
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
				PlayerName:      c.PlayerName,
				RoomID:          playerMessage.Content,
				Conn:            c.Conn,
				ReceiveGameChan: c.ReceiveGameChan,
			}
		case setPlayerName:
			c.PlayerName = string(playerMessage.Content)
			c.Hub.RegisterClientChan <- &RegisterClientRequest{PlayerName: c.PlayerName, Client: c}
			c.sendRoomList()
		case prepare:
			playerMessage.Content = fmt.Sprintf("%d", c.Room.Position)
			message, _ := json.Marshal(playerMessage)
			c.Room.HandlerChan <- message
		default:
			log.Println("用户信息格式错误!")
		}
	}
}

func (c *Client) receiveGameHandler() {
	for game := range c.ReceiveGameChan {
		c.Game = game
		log.Printf("hello %s 新的对局开始了!, 当前打 %s\n", c.PlayerName, game.Round)
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
