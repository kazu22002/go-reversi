package player

type Player struct {
	HumanPlayer bool
	CellType    uint8
}

func New( humanPlayer bool, cellType uint8) Player {
	return Player{
		HumanPlayer: humanPlayer,
		CellType:    cellType,
	}
}
