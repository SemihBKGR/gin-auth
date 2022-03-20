package main

import (
	"gin-auth/handle"
	"github.com/gin-gonic/gin"
)

const serverPortEnv = "GIN_PORT"
const jwtSecretEnv = "GIN_JWT_SECRET"

const serverDefaultPort = 9000
const jwtSecretDefault = "s3cr3t"

const jwtIssuer = "gin-auth"

func routeHandlerFuncs(e *gin.Engine) {

	e.GET("/health",
		handle.Health)

	e.POST("/login",
		handle.Login(loginService, jwtService))

	e.POST("/user",
		handle.SaveUser(userRepo, passEncoder))

	e.PUT("/user/:id",
		handle.JwtAuthenticationMw(jwtService),
		handle.UpdateUser(userRepo, passEncoder))

	e.GET("/user/:id",
		handle.FindUser(userRepo))

	e.DELETE("/user/:id",
		handle.DeleteUser(userRepo))

}
