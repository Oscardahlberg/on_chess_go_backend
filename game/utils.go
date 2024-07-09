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
	NewState    GameData
}

type GameData struct {
	State           GameState
	WhiteLostPieces []string
	BlackLostPieces []string
}

type GameState [8][8]string

func InitChess() GameState {
	return GameState{
		{WTower, WPawn, "", "", "", "", BPawn, BTower},
		{WKnight, WPawn, "", "", "", "", BPawn, BKnight},
		{WBishop, WPawn, "", "", "", "", BPawn, BBishop},
		{WQueen, WPawn, "", "", "", "", BPawn, BQueen},
		{WKing, WPawn, "", "", "", "", BPawn, BKing},
		{WBishop, WPawn, "", "", "", "", BPawn, BBishop},
		{WKnight, WPawn, "", "", "", "", BPawn, BKnight},
		{WTower, WPawn, "", "", "", "", BPawn, BTower},
	}
}
