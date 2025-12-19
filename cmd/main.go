package main

import (
	"net/http"
	"time"

	routeV1 "github.com/amitshekhariitbhu/go-backend-clean-architecture/api/route/v1"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/bootstrap"
	"github.com/gin-gonic/gin"
)

func main() {

	app := bootstrap.App()

	env := app.Env

	db := app.Mongo.Database(env.DBName)
	defer app.CloseDBConnection()

	timeout := time.Duration(env.ContextTimeout) * time.Second

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	routerV1 := router.Group("v1")

	routeV1.Setup(env, timeout, db, routerV1)

	router.Run(env.ServerAddress)
}
