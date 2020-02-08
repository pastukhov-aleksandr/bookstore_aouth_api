package app

import (
	"github.com/gin-gonic/gin"
	"github.com/pastukhov-aleksandr/bookstore_aouth_api/src/http"
	"github.com/pastukhov-aleksandr/bookstore_aouth_api/src/repositore/db"
	"github.com/pastukhov-aleksandr/bookstore_aouth_api/src/repositore/rest"
	"github.com/pastukhov-aleksandr/bookstore_aouth_api/src/services/access_token"
)

var router = gin.Default()

func StartApplication() {
	atHandler := http.NewAccessTokenHandler(
		access_token.NewService(rest.NewRestUsersRepository(), db.NewRepository()))

	router.GET("/oauth/access_token/:access_token_id", atHandler.GetById)
	router.POST("/oauth/access_token", atHandler.Create)

	router.Run(":8080")
}
