package setgame

func MakeSet(a int, b int, c int) bool {
	if a < 3 && b < 3 && c < 3 {
		return (a == b && a == c) || (a != b && a != c && b != c)
	}
	return MakeSet(a%3, b%3, c%3) && MakeSet(a/3, b/3, c/3)
}

func getMatch(a int, b int) int {
	if a < 3 && b < 3 {
		return (6 - a - b) % 3
	}
	return getMatch(a%3, b%3) + getMatch(a/3, b/3)*3
}

func HasSet(cards []int) bool {
	for i := 0; i < len(cards); i++ {
		for j := 0; j < i; j++ {
			expected := getMatch(cards[i], cards[j])
			for k := 0; k < j; k++ {
				if cards[k] == expected {
					return true
				}
			}
		}
	}
	return false
}

func PositionsMakeSet(cards []int, positions []int) bool {
	if max(positions) >= len(cards) {
		return false
	}
	return MakeSet(cards[positions[0]], cards[positions[1]], cards[positions[2]])
}

func RemoveAndNormalize(table []int, gone []int, deck []int, positions []int) ([]int, []int, []int, bool) {
	gone = appendByPosition(gone, table, positions)
	table = removeByPosition(table, positions)
	return Normalize(table, gone, deck)
}

func Normalize(table []int, gone []int, deck []int) ([]int, []int, []int, bool) {
	for (len(table) < 12 || !HasSet(table)) && len(deck) > 0 {
		table, deck = showMoreFromDeck(table, deck)
	}
	return table, gone, deck, !HasSet(table)
}

func showMoreFromDeck(table []int, deck []int) ([]int, []int) {
	if len(deck) == 0 {
		return table, deck
	}
	return append(table, deck[len(deck)-3:]...), deck[:len(deck)-3]
}
