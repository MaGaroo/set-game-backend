package messages

type AddCard struct {
	Card   int8
	Row    int8
	Column int8
}

type RemoveCardByCard struct {
	Card int8
}

type RemoveCardByPlace struct {
	Row    int8
	Column int8
}

type ReplaceCard struct {
	Card1 int8
	Card2 int8
}
