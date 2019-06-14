package messages

import "fmt"

type UpdateScore struct {
	Username string
	Score    int8
}

func (m *UpdateScore) ToString() string {
	return fmt.Sprintf("score %s %d", m.Username, m.Score)
}
