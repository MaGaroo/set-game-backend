package messages

type Guess struct {
	PlayerToken string	`json:"token"`
	GuessString string	`json:"guess"`
}
