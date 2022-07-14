package server

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Player struct {
	Name            string
	UUID            string
	Conn            *websocket.Conn
	Prepare         bool
	ReceiveGameChan chan *Game
}

type NewJoinMsg struct {
	Position int32  `json:"position"`
	Name     string `json:"name"`
}

type ExistPlayerMsg struct {
	Name    string `json:"name"`
	Prepare bool   `json:"prepare"`
}

var TeamMateMap = map[int]int{
	0: 1,
	1: 0,
	2: 3,
	3: 2,
}
var OpponentMap = map[int][2]int{
	0: {2, 3},
	1: {3, 2},
	2: {1, 0},
	3: {0, 1},
}

func (p *Player) notifyExistPlayers(existPlayersMsg []*ExistPlayerMsg, idx int) {
	log.Printf("%s notify %+v\n", p.Name, existPlayersMsg)
	sendPlayerMsg := []*ExistPlayerMsg{
		existPlayersMsg[idx],
		nil,
		nil,
		nil,
	}
	teamMateIdx := TeamMateMap[idx]
	if teamMateIdx+1 <= len(existPlayersMsg) {
		sendPlayerMsg[1] = existPlayersMsg[teamMateIdx]
	}
	oppose1Idx, oppose2Idx := OpponentMap[idx][0], OpponentMap[idx][1]
	if oppose1Idx+1 <= len(existPlayersMsg) {
		sendPlayerMsg[2] = existPlayersMsg[oppose1Idx]
	}
	if oppose2Idx+1 <= len(existPlayersMsg) {
		sendPlayerMsg[3] = existPlayersMsg[oppose2Idx]
	}
	content, _ := json.Marshal(sendPlayerMsg)
	message := ResponseMessage{
		MessageType: existPlayers,
		Content:     string(content),
	}
	sendMessage, _ := json.Marshal(message)
	p.Conn.WriteMessage(websocket.TextMessage, sendMessage)
}
