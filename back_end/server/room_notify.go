package server

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Player struct {
	Name            string
	Conn            *websocket.Conn
	Prepare         bool
	ReceiveGameChan chan *Game
}

type NewJoinMsg struct {
	Position int32  `json:"position"`
	Name     string `json:"name"`
}

type ExistPlayerMsg struct {
	Position int32  `json:"position"`
	Name     string `json:"name"`
	Prepare  bool   `json:"prepare"`
}

func (p *Player) notifyNewPlayerJoin(position int32, name string) {
	content, _ := json.Marshal(NewJoinMsg{
		Position: position,
		Name:     name,
	})
	message := ResponseMessage{
		MessageType: newPlayerJoinRoom,
		Content:     string(content),
	}
	sendMessage, _ := json.Marshal(message)
	p.Conn.WriteMessage(websocket.TextMessage, sendMessage)
}

func (p *Player) notifyExistPlayers(existPlayersMsg []*ExistPlayerMsg) {
	if len(existPlayersMsg) == 0 {
		return
	}
	content, _ := json.Marshal(existPlayersMsg)
	message := ResponseMessage{
		MessageType: existPlayers,
		Content:     string(content),
	}
	sendMessage, _ := json.Marshal(message)
	p.Conn.WriteMessage(websocket.TextMessage, sendMessage)
}

func (p *Player) notifyPlayerPrepare(position int32) {
	message := ResponseMessage{
		MessageType: hasPlayerPrepare,
		Content:     fmt.Sprintf("%d", position),
	}
	sendMessage, _ := json.Marshal(message)
	log.Printf("%s notify prepare %s\n", p.Name, string(sendMessage))
	p.Conn.WriteMessage(websocket.TextMessage, sendMessage)
}
