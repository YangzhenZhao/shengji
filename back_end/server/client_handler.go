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
		log.Println(string(message))
		playerMessage := PlayerMessage{}
		err = json.Unmarshal(message, &playerMessage)
		if err != nil {
			log.Printf("unmarshal playerMessage err: %+v", err)
			continue
		}
		switch playerMessage.MessageType {
		case joinRoom:
			// roomID := playerMessage.Content
		case setPlayerName:
			c.PlayerName = string(playerMessage.Content)
			c.Hub.RegisterClientChan <- &RegisterClientRequest{PlayerName: c.PlayerName, Client: c}
		default:
			log.Println("用户信息格式错误!")
		}
		log.Printf("%+v", playerMessage)
	}
}

func (c *Client) serverMessageHandler() {
	defer func() {
		c.Conn.Close()
	}()
	for {
		select {
		case message := <-c.ReceiveChan:
			serverMessage := &ServerMessage{}
			json.Unmarshal(message, &serverMessage)
			switch serverMessage.MessageType {
			case roomList:
				w, _ := c.Conn.NextWriter(websocket.TextMessage)
				w.Write(message)
				w.Close()
			}
		}
	}
}

func (c *Client) tickerHandler() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
