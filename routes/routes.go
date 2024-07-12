package routes

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"gobackend/handlers"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:3000"
		},
		MaxAge: 12 * time.Hour,
	}))

	// Routes
	router.POST("/new-lobby", handlers.PostNewLobby)
	router.POST("/join-lobby", handlers.PostJoinLobby)
	//router.POST("/create-dumb-bot", bots.PostCreateDumbBot)

	router.GET("/get-lobbys", handlers.GetLobbys)
	router.GET("/", handlers.GetHome)

	// Websockets
	router.GET("/game-channel", func(c *gin.Context) {
		handlers.HandleConnections(c)
	})

	// Temporary Routes
	router.GET("/delete-op-test", handlers.MongoDeleteAllOpenLobbys)
	router.GET("/delete-on-test", handlers.MongoDeleteAllOngoingLobbys)

	return router
}
