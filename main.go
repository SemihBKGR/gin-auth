package main

import (
	"fmt"
	"gin-auth/auth"
	"gin-auth/persist"
	"gin-auth/util"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

var userRepo = persist.NewUserSqliteRepository()
var postRepo = persist.NewPostSqliteRepository()
var commentRepo = persist.NewCommentSqliteRepository()

var passEncoder = auth.NewBcryptPasswordEncoder()

func main() {
	r := gin.Default()
	port := util.GetIntEnvVar(serverPortEnv, serverDefaultPort)
	routeHandlerFuncs(r)
	err := r.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Error(err)
	}
}
