package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/gin-gonic/gin"
	"gobackend/game"
)

var upgrader = websocket.Upgrader{}

type clientMessage struct {
	Server     string
	GameUpdate game.GameUpdateMessage
}

func HandleConnections(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Error upgrading websocket connection: ", err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	defer ws.Close()

	initGame(ws, c)

	for {
		var msg clientMessage

		// From client
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Println("Error upgrading websocket connection: ", err.Error())
			c.AbortWithError(http.StatusInternalServerError, err)
		}

		// Gets the open lobby JSON from Mongodb
		lobby, plyr, err := getOngoingLobby(c)
		if err != nil {
			log.Println("Couldnt find open lobby connected with user id", err.Error())
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		response, err := handleClientMsg(msg, lobby, plyr, ws, c)
		if err != nil {
			log.Println("Client message error: ", err.Error())
			c.AbortWithError(http.StatusInternalServerError, err)
		}

		opponent := lobby.Player1Socket
		if plyr == "Player1" {
			opponent = lobby.Player2Socket
		}

		err = opponent.WriteJSON(response)
		if err != nil {
			log.Println("Error upgrading websocket connection: ", err.Error())
			c.AbortWithError(http.StatusInternalServerError, err)
		}
	}
}

func handleClientMsg(msg clientMessage, lobby Lobby, plyr string, ws *websocket.Conn, c *gin.Context) (clientMessage, error) {
	var err error
	serverMsg := ""
	var gameMsg string
	lobby.GameState = msg.GameUpdate.NewState

	switch msg.GameUpdate.GameMessage {
	case "Game Update":
		gameMsg = "Player1 Turn"
		if plyr == "Player2" {
			gameMsg = "Player2 Turn"
		}
		err = postOngoingLobby(lobby, plyr, c)
	case "White win", "Black win":
		gameMsg = msg.GameUpdate.GameMessage
		err = postEndGame(lobby, plyr, c)
	}

	response := clientMessage{
		Server: serverMsg,
		GameUpdate: game.GameUpdateMessage{
			GameMessage: gameMsg,
			NewState:    msg.GameUpdate.NewState,
		},
	}
	return response, err
}

func initGame(ws *websocket.Conn, c *gin.Context) {
	// Lobby data, Plyr is if the client is host or player2
	lobby, plyr, err := getOngoingLobby(c)
	if err != nil {
		log.Println("Couldnt find open lobby connected with user id", err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	if plyr == "Player1" {
		lobby.Player1Socket = ws
	} else {
		lobby.Player2Socket = ws
	}

	err = postOngoingLobby(lobby, plyr, c) // Make
	if err != nil {
		log.Println("DB error: ", err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	if plyr == "Player2" {
		gameMsg := "false"
		// If the player starting is the one who sent this msg
		if lobby.StartPlayer == plyr {
			gameMsg = "true"
		}

		startMsg := clientMessage{
			Server: plyr,
			GameUpdate: game.GameUpdateMessage{
				GameMessage: gameMsg,
				NewState:    lobby.GameState,
			},
		}

		err = ws.WriteJSON(startMsg)
		if err != nil {
			log.Println("Error writing Json to : ", err.Error())
			c.AbortWithError(http.StatusInternalServerError, err)
		}
	}
}
