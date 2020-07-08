package app

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pastukhov-aleksandr/bookstore_aouth_api/http"
	"github.com/pastukhov-aleksandr/bookstore_aouth_api/repositore/db"
	"github.com/pastukhov-aleksandr/bookstore_aouth_api/repositore/rest"
	"github.com/pastukhov-aleksandr/bookstore_aouth_api/services/access_token"
)

var router = gin.Default()

func StartApplication() {
	atHandler := http.NewAccessTokenHandler(
		access_token.NewService(rest.NewRestUsersRepository(), db.NewRepository()))

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:8081"}
	router.Use(cors.New(config))

	router.GET("/oauth/access_token/:access_token_id", atHandler.GetById)
	router.POST("/oauth/access_token", atHandler.Create)
	router.POST("/oauth/refresh_token", atHandler.Create)
	router.POST("/oauth/logout", atHandler.DeleteRefreshToken)

	router.Run(":8080")
}
