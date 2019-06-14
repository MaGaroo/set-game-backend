package messages

import "fmt"

type UpdateCard struct {
	Card   int
	Row    int
	Column int
}

func (m UpdateCard) ToBytes() []byte {
	return []byte(fmt.Sprintf("card %d %d %d", m.Card, m.Row, m.Column))
}
