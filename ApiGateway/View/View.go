package View

import (
	"github.com/gin-gonic/gin"
	"github.com/mhthrh/BlueBank/ApiGateway/Controller"
	"github.com/mhthrh/BlueBank/Pool"
	"net/http"
)

func New(c *chan Pool.Connection) {
	Controller.New(c)
}
func RunSync() http.Handler {
	router := gin.New()
	//router.Use(Controller.Middleware)
	router.Use(gin.Recovery())

	router.POST("/gateway/login", Controller.GatewaySignIn)
	router.POST("/user/signup", Controller.UserSignUp)
	router.POST("/user/signin", Controller.UserSignIn)
	//router.GET("/version", Controller.Version)
	//
	//router.NoRoute(Controller.NotFound)

	return router
}
func RunAsync() http.Handler {
	router := gin.New()
	router.Use(gin.Recovery())
	//router.GET("/user/message", Controller.Message)
	//
	//router.NoRoute(Controller.NotFound)

	return router
}
