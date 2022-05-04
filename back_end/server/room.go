package server

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/gorilla/websocket"
)

const maxRoomPlayer = 4

type Room struct {
	ID      string
	Players []*Player

	ClientHandlerChan chan []byte
}

func (r *Room) RegisterPlayer(joinRoomRequest *JoinRoomRequest) (bool, int, chan []byte, string) {
	for _, player := range r.Players {
		if player.Name == joinRoomRequest.PlayerName {
			return false, -1, nil, "不可重复加入!"
		}
	}
	if len(r.Players) == maxRoomPlayer {
		return false, -1, nil, "房间已满!"
	}

	newPlayer := &Player{
		joinRoomRequest.PlayerName,
		joinRoomRequest.Conn,
		false,
		joinRoomRequest.ReceiveGameChan,
	}
	r.Players = append(r.Players, newPlayer)

	var existPlayers []*ExistPlayerMsg
	for i, player := range r.Players {
		existPlayers = append(existPlayers, &ExistPlayerMsg{
			Position: int32(i + 1),
			Name:     player.Name,
			Prepare:  player.Prepare,
		})
	}

	for _, player := range r.Players {
		player.notifyExistPlayers(existPlayers)
	}

	return true, len(r.Players), r.ClientHandlerChan, ""
}

func (r *Room) clientMessageHandler() {
	for {
		message := <-r.ClientHandlerChan
		playerMessage := RequestMessage{}
		err := json.Unmarshal(message, &playerMessage)
		if err != nil {
			log.Printf("unmarshal playerMessage err: %+v", err)
			continue
		}
		switch playerMessage.MessageType {
		case prepare:
			position, _ := strconv.ParseInt(playerMessage.Content, 10, 32)
			r.Players[position-1].Prepare = true
			for i, player := range r.Players {
				if i+1 == int(position) {
					continue
				}
				player.notifyPlayerPrepare(int32(position))
			}
			if r.isAllPrepare() {
				match := &Match{
					FirstTeamRound:   "2",
					SecondTeamRound:  "2",
					PlayerConns:      r.playerConns(),
					ReceiveGameChans: r.receiveGameChans(),
				}
				go match.Run()
			}
		default:
			log.Println("用户信息格式错误!")
		}
	}
}

func (r *Room) isAllPrepare() bool {
	if len(r.Players) != 4 {
		return false
	}
	for _, player := range r.Players {
		if !player.Prepare {
			return false
		}
	}
	return true
}

func (r *Room) playerConns() []*websocket.Conn {
	var res []*websocket.Conn
	for _, player := range r.Players {
		res = append(res, player.Conn)
	}
	return res
}

func (r *Room) receiveGameChans() []chan *Game {
	var res []chan *Game
	for _, player := range r.Players {
		res = append(res, player.ReceiveGameChan)
	}
	return res
}
