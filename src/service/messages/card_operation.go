package messages

import "fmt"

type UpdateCard struct {
	Card   int8
	Row    int8
	Column int8
}

func (m UpdateCard) ToString() string {
	return fmt.Sprintf("card %d %d %d", m.Card, m.Row, m.Column)
}
