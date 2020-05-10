package api

import (
	"bytes"
	"encoding/binary"
)

type Position struct {
	X uint32
	Y uint32
}

func (p *Position) GetBytes() []byte {
	buf := bytes.NewBuffer(make([]byte, 0, 8))
	binary.Write(buf, binary.BigEndian, p.X)
	binary.Write(buf, binary.BigEndian, p.Y)

	return buf.Bytes()
}

type GameStats struct {
	PlayerID   string
	NumPlayers int
	Mass       float64
	FoodEaten  int
	CellsEaten int
	Rank       int
}

type MouseEvent struct {
	SeqID uint32
	X     float32
	Y     float32
}
