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
		handle.Health,
	)

	e.POST("/login",
		handle.Login(loginService, jwtService),
	)

	e.POST("/user",
		handle.SaveUser(userRepo, passEncoder),
	)

	e.PUT("/user",
		handle.JwtAuthenticationMw(jwtService),
		handle.UpdateUser(userRepo, passEncoder),
	)

	e.GET("/user",
		handle.JwtAuthenticationMw(jwtService),
		handle.FindUser(userRepo),
	)

	e.GET("/user/:id",
		handle.JwtAuthenticationMw(jwtService),
		handle.FindUserById(userRepo),
	)

	e.DELETE("/user/:id",
		handle.JwtAuthenticationMw(jwtService),
		handle.JwtAuthorizationHasEachRoleMv("ADMIN"),
		handle.DeleteUser(userRepo),
	)

	e.POST("/post",
		handle.JwtAuthenticationMw(jwtService),
		handle.SavePost(postRepo),
	)

}
