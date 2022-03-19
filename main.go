package main

import (
	"fmt"
	"gin-auth/persist"
	"gin-auth/util"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func main() {
	r := gin.Default()
	port := util.GetIntEnvVar(serverPortEnv, serverDefaultPort)
	err := r.Run(fmt.Sprintf(":%d", port))
	_ = persist.NewUserSqliteRepository()
	if err != nil {
		log.Error(err)
	}
}
