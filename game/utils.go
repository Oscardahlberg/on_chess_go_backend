package game

type piece struct {
	Name           string
	Colour         string
	Move           string
	AdditionalMove string
}

type GameUpdateMessage struct {
	GameMessage string
	NewState    GameData
}

type GameData struct {
	State           GameState
	WhiteLostPieces []string
	BlackLostPieces []string
}

type GameState [8][8]*piece

func InitChess() GameState {
	return GameState{
		{btower, bknight, bbishop, bking, bqueen, bbishop, bknight, btower},
		{bpawn, bpawn, bpawn, bpawn, bpawn, bpawn, bpawn, bpawn},
		{nil, nil, nil, nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil, nil, nil, nil},
		{wpawn, wpawn, wpawn, wpawn, wpawn, wpawn, wpawn, wpawn},
		{wtower, wknight, wbishop, wking, wqueen, wbishop, wknight, wtower},
	}
}

var wpawn = &piece{
	Name:           "pawn",
	Colour:         "black",
	Move:           "1;Vertical",
	AdditionalMove: "1;Diagonal;UP;A",
}

var bpawn = &piece{
	Name:           "pawn",
	Colour:         "black",
	Move:           "1;Vertical",
	AdditionalMove: "1;Diagonal;UP;A",
}

var wtower = &piece{
	Name:           "tower",
	Colour:         "white",
	Move:           "8;Vertical;Horizontal;",
	AdditionalMove: "",
}

var btower = &piece{
	Name:           "tower",
	Colour:         "black",
	Move:           "8;Vertical;Horizontal;",
	AdditionalMove: "",
}

var wknight = &piece{
	Name:           "knight",
	Colour:         "white",
	Move:           "",
	AdditionalMove: "5;Knight;", // knight move
}

var bknight = &piece{
	Name:           "knight",
	Colour:         "black",
	Move:           "",
	AdditionalMove: "5;Knight;", // knight move
}

var wbishop = &piece{
	Name:           "bishop",
	Colour:         "white",
	Move:           "8;Diagonal;",
	AdditionalMove: "",
}

var bbishop = &piece{
	Name:           "bishop",
	Colour:         "black",
	Move:           "8;Diagonal;",
	AdditionalMove: "",
}

var wqueen = &piece{
	Name:           "queen",
	Colour:         "white",
	Move:           "8;Any",
	AdditionalMove: "",
}

var bqueen = &piece{
	Name:           "queen",
	Colour:         "black",
	Move:           "8;Any",
	AdditionalMove: "",
}

var wking = &piece{
	Name:           "king",
	Colour:         "white",
	Move:           "1;Any",
	AdditionalMove: "King",
}

var bking = &piece{
	Name:           "king",
	Colour:         "black",
	Move:           "1;Any",
	AdditionalMove: "King",
}

/*
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
*/
