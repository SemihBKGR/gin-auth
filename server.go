package main

import (
	"gin-auth/handle"
	"github.com/gin-gonic/gin"
)

const serverPortEnv = "GIN_PORT"

const serverDefaultPort = 9000

func routeHandlerFuncs(e *gin.Engine) {
	e.GET("/health", handle.Health)
	e.POST("/user", handle.SaveUser(userRepo, passEncoder))
	e.PUT("/user/:id", handle.UpdateUser(userRepo, passEncoder))
	e.GET("/user/:id", handle.FindUser(userRepo))
	e.DELETE("/user/:id", handle.DeleteUser(userRepo))
}
