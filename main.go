package main

import (
	"fmt"
	"gin-auth/util"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func main() {
	r := gin.Default()
	port := util.GetIntEnvVar(serverPortEnv, serverDefaultPort)
	err := r.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Error(err)
	}
}
