package main

import (
	"fmt"
	"gin-auth/auth"
	"gin-auth/auth/jwt"
	"gin-auth/persist"
	"gin-auth/util"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var log = logrus.New()

var userRepo = persist.NewUserSqliteRepository()
var postRepo = persist.NewPostSqliteRepository()
var commentRepo = persist.NewCommentSqliteRepository()

var passEncoder = auth.NewBcryptPasswordEncoder()
var loginService = auth.NewDefaultLoginService(userRepo, passEncoder)

var jwtService = jwt.NewJwtService(util.GetEnvVar(jwtSecretEnv, jwtSecretDefault), jwtIssuer)

func init() {
	persist.InitDatabase(func(db *gorm.DB) {
		db.Exec(auth.InsertRolesQuery)
	})
}

func main() {
	r := gin.Default()
	port := util.GetIntEnvVar(serverPortEnv, serverDefaultPort)
	routeHandlerFuncs(r)
	err := r.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Error(err)
	}
}
