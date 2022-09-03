package dto

type MasterLevel int

const (
	NormalLevelMaster MasterLevel = iota
	PlayNumberLevelMaster
	JokerLevelMaster
)

// GetCardMasterLevel ...
func GetCardMasterLevel(card *Poker, playNumber string) MasterLevel {
	if card.Number == "joker" {
		return JokerLevelMaster
	}
	if card.Number == playNumber {
		return PlayNumberLevelMaster
	}
	return NormalLevelMaster
}
