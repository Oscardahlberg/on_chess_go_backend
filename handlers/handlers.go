package handlers

import (
	"log"
	"math/rand"
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

type Lobby struct {
	Id            string
	Player1       string
	Player1Socket *websocket.Conn
	Player2       string
	Player2Socket *websocket.Conn
	Expires       time.Time
	StartPlayer   string
	GameState     GameState
}

func GetHome(c *gin.Context) {
	c.String(http.StatusOK, "Welcome home oniichan\n")
}

func PostNewLobby(c *gin.Context) {
	userId := getCookie(c)

	del_res, err := database.OpenLobbys.DeleteMany(context.Background(), bson.M{"Player1": userId})
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: ", err))
		return
	}
	if del_res != nil {
		log.Println("Response deleting open lobby: ", del_res)
	}

	del_res, err = database.OngoingLobbys.DeleteMany(context.Background(), bson.M{"Player1": userId})
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: ", err))
		return
	}
	if del_res != nil {
		log.Println("Response deleting ongoing lobby: ", del_res)
	}

	gameId := uuid.NewString()
	expirationTime := time.Now().Add(24 * time.Hour)

	ins_res, err := database.OpenLobbys.InsertOne(context.Background(), Lobby{
		Id:            gameId,
		Player1:       userId,
		Player1Socket: nil,
		Player2:       nil,
		Player2Socket: nil,
		Expires:       expirationTime,
		StartPlayer:   nil,
		GameState:     game.initChess(),
	})
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
		return
	}
	if ins_res != nil {
		log.Println("Response inserting open lobby: ", ins_res)
	}

	c.JSON(http.StatusOK, bson.M{"GameId": gameId})
}

func PostJoinLobby(c *gin.Context) {
	userId := getCookie(c)

	type GameId struct {
		LobbyId string
	}
	var gameId GameId
	if err := c.BindJSON(&gameId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err: ": err.Error()})
		return
	}

	result := database.OpenLobbys.FindOne(context.Background(), bson.M{"Id": gameId.LobbyId})

	var parsedLobby Lobby
	if err := result.Decode(&parsedLobby); err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
		return
	}

	if parsedLobby.Player2 != nil {
		c.String(http.StatusNotFound, "")
		return
	}
	parsedLobby.Player2 = userId

	parsedLobby.StartPlayer = "Player1"
	if rand.Int(2) == 2 {
		parsedLobby.White = "Player2"
	}

	delRes, err := database.OpenLobbys.DeleteOne(context.Background(), bson.M{"Id": gameId.LobbyId})
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
		return
	}
	if delRes.DeletedCount == 0 {
		c.String(
			http.StatusInternalServerError,
			fmt.Sprintf("Open lobby with id %v not found", gameId.LobbyId))
		return
	}

	insRes, err := database.OngoingLobbys.InsertOne(context.Background(), parsedLobby)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
		return
	}
	log.Println("Response deleting open lobby: ", delRes)
	log.Println("Response inserting ongoing lobby: ", insRes)

	c.String(http.StatusOK, "")
}

func GetLobbys(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")

	type lobbys struct {
		Open   []Lobby
		Closed []Lobby
	}

	var parsed_open_lobbys []Lobby
	var parsed_ongoing_lobbys []Lobby

	//todo: check for newest 10
	cursor, err := database.OpenLobbys.Find(
		context.Background(), bson.M{}, options.Find())
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
		return
	}
	defer cursor.Close(context.Background())

	decodeLobby(c, cursor, &parsed_open_lobbys)

	//todo: check for newest 10
	cursor, err = database.OngoingLobbys.Find(
		context.Background(), bson.M{}, options.Find().SetLimit(10))
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
		return
	}
	defer cursor.Close(context.Background())

	decodeLobby(c, cursor, &parsed_ongoing_lobbys)
	c.JSON(http.StatusOK, lobbys{parsed_open_lobbys, parsed_ongoing_lobbys})
}

func getCookie(c *gin.Context) string {
	userId, err := c.Cookie("userId")
	if err != nil {
		userId = uuid.NewString()
		c.SetCookie("userId", userId, 7200, "/", "localhost", false, true)
	}
	return userId
}

func decodeLobby(c *gin.Context, cursor *mongo.Cursor, lobbys *[]Lobby) {
	if err := cursor.All(context.Background(), lobbys); err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Cursor decode error: %v", err))
		return
	}
}

func getOngoingLobby(c *gin.Context) (Bson.M, string, error) {
	userId := getCookie(c)

	var result bson.M
	plyr := 0

	err := database.OngoingLobbys.FindOne(context.Background(), Bson.M{"Player1": userId}).Decode(&result)
	if err == nil {
		return result, "Player1", nil
	}

	err = database.OngoingLobbys.FindOne(context.Background(), Bson.M{"Player2": userId}).Decode(&result)
	if err == nil {
		return result, "Player2", nil
	}

	return nil, nil, errors.New("No open lobby found ")
}

// for testing
func MongoPopulate(c *gin.Context) {
	ins_res, err := database.OpenLobbys.InsertOne(context.Background(), bson.M{
		"id":      uuid.NewString(),
		"host":    uuid.NewString(),
		"player2": "",
		"expires": time.Now().Add(24 * time.Hour),
	})
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
		return
	}
	if ins_res != nil {
		log.Println("Response inserting open lobby: ", ins_res)
	}

	c.String(http.StatusOK, "Inserted a open lobby\n")
}

func MongoDeleteAllOpenLobbys(c *gin.Context) {
	_, err := database.OpenLobbys.DeleteMany(context.Background(), bson.M{})

	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
		return
	}

	c.String(http.StatusOK, "Deleted all Open Lobbys\n")
}

func MongoDeleteAllOngoingLobbys(c *gin.Context) {
	_, err := database.OngoingLobbys.DeleteMany(context.Background(), bson.M{})

	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
		return
	}

	c.String(http.StatusOK, "Deleted all Closed Lobbys\n")
}

// for testing
func MongoGet(c *gin.Context) {
	cursor, err := database.OpenLobbys.Find(
		context.Background(), bson.M{})
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
		return
	}
	defer cursor.Close(context.Background())

	var lobbys []Lobby

	if err := cursor.All(context.Background(), &lobbys); err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Cursor decode error: %v", err))
		return
	}
	log.Println(lobbys)

	if err = cursor.Err(); err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, lobbys)
}
