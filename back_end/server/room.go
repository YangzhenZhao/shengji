package server

import (
	"encoding/json"
	"log"
	"strconv"
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
	var existPlayers []*ExistPlayerMsg
	for i, player := range r.Players {
		player.notifyNewPlayerJoin(int32(len(r.Players))+1, joinRoomRequest.PlayerName)
		existPlayers = append(existPlayers, &ExistPlayerMsg{
			Position: int32(i + 1),
			Name:     player.Name,
			Prepare:  player.Prepare,
		})
	}
	newPlayer := &Player{
		joinRoomRequest.PlayerName,
		joinRoomRequest.Conn,
		false,
	}
	newPlayer.notifyExistPlayers(existPlayers)
	r.Players = append(r.Players, newPlayer)
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
