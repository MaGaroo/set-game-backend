package messages

import (
	"fmt"
)

type EndGame struct {
	Winner string
}

func (m EndGame) ToBytes() []byte {
	return []byte(fmt.Sprintf("the_end %s", m.Winner))
}
