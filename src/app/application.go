package app

import (
	"github.com/gin-gonic/gin"
	"github.com/pastukhov-aleksandr/bookstore_aouth_api/src/domain/access_token"
	"github.com/pastukhov-aleksandr/bookstore_aouth_api/src/http"
	"github.com/pastukhov-aleksandr/bookstore_aouth_api/src/repositore/db"
)

var router = gin.Default()

func StartApplication() {
	atHandler := http.NewAccessTokenHandler(access_token.NewService(db.NewRepository()))

	router.GET("/oauth/access_token/:access_token_id", atHandler.GetById)
	router.POST("/oauth/access_token", atHandler.Create)

	router.Run(":8080")

}
