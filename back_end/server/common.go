package server

func getRelativePos(myPos, referPos int) int {
	if myPos == referPos {
		return 0
	}
	if TeamMateMap[myPos] == referPos {
		return 1
	}
	if OpponentMap[myPos][0] == referPos {
		return 2
	}
	return 3
}
