package main

import (
	"gin-auth/auth"
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

	e.POST("/post",
		handle.JwtAuthenticationMw(jwtService),
		handle.SavePost(postRepo),
	)

	e.PUT("/post/:id",
		handle.JwtAuthenticationMw(jwtService),
		handle.UpdatePost(postRepo),
	)

	e.GET("/post/:id",
		handle.JwtAuthenticationMw(jwtService),
		handle.FindPost(postRepo),
	)

	e.GET("/post/list",
		handle.JwtAuthenticationMw(jwtService),
		handle.FindAllPosts(postRepo),
	)

	e.GET("/post/list/:username",
		handle.JwtAuthenticationMw(jwtService),
		handle.FindAllPostsByUsername(postRepo),
	)

	e.DELETE("/post/:id",
		handle.JwtAuthenticationMw(jwtService),
		handle.DeletePost(postRepo),
	)

	e.DELETE("/post/force/:id",
		handle.JwtAuthenticationMw(jwtService),
		handle.JwtAuthorizationHasAnyRoleMv(auth.RoleAdmin, auth.RoleManager),
		handle.DeletePostForcibly(postRepo),
	)

}
