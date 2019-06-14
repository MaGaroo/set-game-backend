package messages

import (
	"fmt"
	"strings"
)

type IntroRequest struct {
	RoomToken   string
	Username    string
	PlayerToken string
}

func ParseIntroRequest(message string) *IntroRequest {
	args := strings.Fields(message)
	if args[0] != "intro" || (len(args) != 4 && len(args) != 3) {
		return nil
	}
	token := ""
	if len(args) == 4 {
		token = args[3]
	}
	return &IntroRequest{
		RoomToken:   args[1],
		Username:    args[2],
		PlayerToken: token,
	}
}

type IntroResponse struct {
	Username    string
	PlayerToken string
}

func (m IntroResponse) ToString() string {
	return fmt.Sprintf("intro %s %s", m.Username, m.PlayerToken)
}
