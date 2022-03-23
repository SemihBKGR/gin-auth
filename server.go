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

	e.Use(handle.JwtAuthenticationMw(jwtService))

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
		handle.JwtAuthenticationRequiredMw(jwtService),
		handle.UpdateUser(userRepo, passEncoder),
	)

	e.GET("/user",
		handle.JwtAuthenticationRequiredMw(jwtService),
		handle.FindUser(userRepo),
	)

	e.GET("/user/:username",
		handle.JwtAuthenticationRequiredMw(jwtService),
		handle.FindUserByUsername(userRepo),
	)

	e.POST("/post",
		handle.JwtAuthenticationRequiredMw(jwtService),
		handle.SavePost(postRepo),
	)

	e.PUT("/post/:id",
		handle.JwtAuthenticationRequiredMw(jwtService),
		handle.UpdatePost(postRepo),
	)

	e.PUT("/post/force/:id",
		handle.JwtAuthenticationRequiredMw(jwtService),
		handle.JwtAuthorizationHasAnyRoleMv(auth.RoleAdmin, auth.RoleManager),
		handle.UpdatePostForcibly(postRepo),
	)

	e.GET("/post/:id",
		handle.JwtAuthenticationRequiredMw(jwtService),
		handle.FindPost(postRepo),
	)

	e.GET("/post/list",
		handle.JwtAuthenticationRequiredMw(jwtService),
		handle.FindAllPosts(postRepo),
	)

	e.GET("/post/list/:username",
		handle.JwtAuthenticationRequiredMw(jwtService),
		handle.FindAllPostsByUsername(postRepo),
	)

	e.DELETE("/post/:id",
		handle.JwtAuthenticationRequiredMw(jwtService),
		handle.DeletePost(postRepo),
	)

	e.DELETE("/post/force/:id",
		handle.JwtAuthenticationRequiredMw(jwtService),
		handle.JwtAuthorizationHasAnyRoleMv(auth.RoleAdmin, auth.RoleManager),
		handle.DeletePostForcibly(postRepo),
	)

	e.POST("/comment/:postId",
		handle.JwtAuthenticationRequiredMw(jwtService),
		handle.SaveComment(commentRepo),
	)

	e.PUT("/comment/:id",
		handle.JwtAuthenticationRequiredMw(jwtService),
		handle.UpdateComment(commentRepo),
	)

	e.PUT("/comment/force/:id",
		handle.JwtAuthenticationRequiredMw(jwtService),
		handle.JwtAuthorizationHasAnyRoleMv(auth.RoleAdmin, auth.RoleManager, auth.RoleModerator),
		handle.UpdateCommentForcibly(commentRepo),
	)

	e.GET("/comment/list",
		handle.JwtAuthenticationRequiredMw(jwtService),
		handle.FindAllComments(commentRepo),
	)

	e.GET("/comment/list/:username",
		handle.JwtAuthenticationRequiredMw(jwtService),
		handle.FindAllCommentsByUsername(commentRepo),
	)

	e.DELETE("/comment/:id",
		handle.JwtAuthenticationRequiredMw(jwtService),
		handle.DeleteComment(commentRepo),
	)

	e.DELETE("/comment/force/:id",
		handle.JwtAuthenticationRequiredMw(jwtService),
		handle.JwtAuthorizationHasAnyRoleMv(auth.RoleAdmin, auth.RoleManager, auth.RoleModerator),
		handle.DeleteCommentForcibly(commentRepo),
	)

	e.PUT("/role/:username",
		handle.JwtAuthenticationRequiredMw(jwtService),
		handle.JwtAuthorizationHasEachRoleMv(auth.RoleAdmin),
		handle.AddRole(userRepo),
	)

	e.DELETE("/role/:username",
		handle.JwtAuthenticationRequiredMw(jwtService),
		handle.JwtAuthorizationHasEachRoleMv(auth.RoleAdmin),
		handle.RemoveRole(userRepo),
	)

}
