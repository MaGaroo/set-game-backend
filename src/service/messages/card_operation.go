package messages

import (
	"fmt"
	"encoding/json"
)

type UpdateCard struct {
	Cards []int
}

func (m UpdateCard) ToBytes() []byte {
	if cards, err := json.Marshal(m.Cards); err == nil {
		return []byte(fmt.Sprintf("card %s", cards))
	}
	return []byte{}
}
