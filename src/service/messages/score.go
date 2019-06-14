package messages

import "fmt"

type UpdateScore struct {
	Username string
	Score    int8
}

func (m *UpdateScore) ToBytes() []byte {
	return []byte(fmt.Sprintf("score %s %d", m.Username, m.Score))
}
