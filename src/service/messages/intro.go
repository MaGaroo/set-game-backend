package messages

import "fmt"

type IntroRequest struct {
	GameToken   string
	Username    string
	PlayerToken string
}

type IntroResponse struct {
	Username    string
	PlayerToken string
}

func (m IntroResponse) ToString() string {
	return fmt.Sprintf("intro %s %s", m.Username, m.PlayerToken)
}
