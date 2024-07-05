package game

const (
	WPawn   = "WPAWN"
	WTower  = "WTOWER"
	WKnight = "WKNIGHT"
	WBishop = "WBISHOP"
	WQueen  = "WQUEEN"
	WKing   = "WKING"

	BPawn   = "BPAWN"
	BTower  = "BTOWER"
	BKnight = "BKNIGHT"
	BBishop = "BBISHOP"
	BQueen  = "BQUEEN"
	BKing   = "BKING"
)

type GameUpdateMessage struct {
	GameMessage string
	NewState    GameState
}

var GameState [8][8]string

func initChess() GameState {
	GameState{
		{WTower, WPawn, nil, nil, nil, nil, BPawn, BTower},
		{WKnight, WPawn, nil, nil, nil, nil, BPawn, BKnight},
		{WBishop, WPawn, nil, nil, nil, nil, BPawn, BBishop},
		{WQueen, WPawn, nil, nil, nil, nil, BPawn, BQueen},
		{WKing, WPawn, nil, nil, nil, nil, BPawn, BKing},
		{WBishop, WPawn, nil, nil, nil, nil, BPawn, BBishop},
		{WKnight, WPawn, nil, nil, nil, nil, BPawn, BKnight},
		{WTower, WPawn, nil, nil, nil, nil, BPawn, BTower},
	}
}
