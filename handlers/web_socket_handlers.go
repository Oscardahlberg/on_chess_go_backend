package handlers

import (
	"log"
	"net/http"

	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gobackend/database"
	"gobackend/game"
)

var upgrader = websocket.Upgrader{}

func handleConnections(c *gin.Context) {
	type clientMessage struct {
		Server     string
		GameUpdate GameUpdateMessage
	}

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Error upgrading websocket connection: ", err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	defer ws.Close()

	// Sends the start message to whoever player starts
	initGame(ws, &c)

	for {
		var msg clientMessage

		// From client
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Println("Error upgrading websocket connection: ", err.Error())
			c.AbortWithError(http.StatusInternalServerError, err)
		}

		// Gets the open lobby JSON from Mongodb
		lobby, Plyr, err := routes.getOngoingLobby(&c)
		if err != nil {
			log.Println("Couldnt find open lobby connected with user id", err.Error())
			c.AbortWithError(http.StatusInternalServerError, err)
		}

		response := gameHandler.handleMsg(Plyr, msg, lobby)
		if err != nil {
			log.Println("Error with message from client ", err.Error())
			c.AbortWithError(http.StatusInternalServerError, err)
		}

		// To client
		err = ws.WriteJson(response)
		if err != nil {
			log.Println("Error upgrading websocket connection: ", err.Error())
			c.AbortWithError(http.StatusInternalServerError, err)
		}
	}
}

// Since player 2 will always be the second player joining the lobby
// He will have to send the game start message to player 1
// Firstly the websocket connection info gets stored
func initGame(ws *websocket.Conn, c *gin.Context) {
	lobby, Plyr, err := routes.getOngoingLobby(&c)
	if err != nil {
		log.Println("Couldnt find open lobby connected with user id", err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	if Plyr == "Player1" {
		lobby.Player1Socket = ws
	}
	lobby.Player2Socket = ws

	err = routes.postOngoingLobby() // Make
	if err != nil {
		log.Println("DB error: ", err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	if Plyr == "Player2" {
		gameMsg := "Player1 Start"
		if lobby.StartPlayer == "Player2" {
		gameMsg:
			"Player2 Start"
		}

		startMsg = clientMessage{
			Server: "Player 2",
			GameUpdate: GameUpdateMessage{
				GameMessage: gameMsg,
				NewState:    lobby.GameState, // Will always be start state
			},
		}

		err = lobby.Player1Socket.WriteJson(startMsg)
		if err != nil {
			log.Println("Error writing Json to : ", err.Error())
			c.AbortWithError(http.StatusInternalServerError, err)
		}
	}
}
