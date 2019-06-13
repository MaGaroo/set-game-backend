package messages

type IntroRequest struct {
	GameToken   string
	Username    string
	PlayerToken string
}

type IntroResponse struct {
	PlayerToken string
}
