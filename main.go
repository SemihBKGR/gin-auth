package main

import (
	"fmt"
	"gin-auth/handle"
	"gin-auth/persist"
	"gin-auth/util"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func main() {
	r := gin.Default()
	port := util.GetIntEnvVar(serverPortEnv, serverDefaultPort)
	userRepo := persist.NewUserSqliteRepository()
	_ = persist.NewPostSqliteRepository()
	_ = persist.NewCommentSqliteRepository()
	r.GET("/health", handle.Health)
	r.POST("/user", handle.SaveUser(userRepo))
	r.GET("/user/:id", handle.FindUser(userRepo))
	err := r.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Error(err)
	}
}
