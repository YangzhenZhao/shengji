package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/YangzhenZhao/shengji/back_end/server/dto"
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
		log.Printf("player message: %+v\n", playerMessage)
		switch playerMessage.MessageType {
		case joinRoom:
			c.Hub.JoinRoomRequestChan <- &JoinRoomRequest{
				PlayerUUID:      c.UUID,
				PlayerName:      c.PlayerName,
				RoomID:          playerMessage.Content,
				Conn:            c.Conn,
				ReceiveGameChan: c.ReceiveGameChan,
			}
		case setPlayerName:
			playerMeta, err := dto.GetPlayerMeta(playerMessage.Content)
			if err != nil {
				log.Printf("unmarshal playerMeta err: %+v", err)
				continue
			}
			c.PlayerName = playerMeta.PlayerName
			c.UUID = playerMeta.UUID
			c.Hub.RegisterClientChan <- &RegisterClientRequest{PlayerName: c.PlayerName, PlayerUUID: c.UUID, Client: c}
			c.sendRoomList()
		case prepare:
			playerMessage.Content = fmt.Sprintf("%d", c.Room.Position)
			message, _ := json.Marshal(playerMessage)
			c.Room.HandlerChan <- message
		case showMaster:
			var req dto.ShowMasterRequest
			err := json.Unmarshal([]byte(playerMessage.Content), &req)
			if err != nil {
				log.Printf("unmarshal playerMeta err: %+v", err)
				continue
			}
			fmt.Printf("%+v\n", req)
			c.Game.receiveShowMaster(&dto.GameShowMasterRequest{
				Req:      req,
				Position: c.Room.Position,
			})
		case showMasterDone:
			c.Game.receiveShowMasterDone()
		case kouCards:
			var kouCards []*dto.Poker
			err := json.Unmarshal([]byte(playerMessage.Content), &kouCards)
			if err != nil {
				log.Printf("unmarshal kouCards err: %+v", err)
				continue
			}
			c.Game.receiveBottomCards(kouCards)
		case playCards:
			var playCards []*dto.Poker
			err := json.Unmarshal([]byte(playerMessage.Content), &playCards)
			if err != nil {
				log.Printf("unmarshal playCards err: %+v", err)
				continue
			}
			c.Game.receivePlayCards(int(c.Room.Position)-1, playCards)
		default:
			log.Println("用户信息格式错误!")
		}
	}
}

func (c *Client) receiveGameHandler() {
	for game := range c.ReceiveGameChan {
		c.Game = game
		log.Printf("hello %s 新的对局开始了!, 当前打 %s %+v\n", c.PlayerName, game.Round, game)
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
