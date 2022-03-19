package handle

import (
	"encoding/json"
	"gin-auth/persist"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
)

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, struct{ Status string }{Status: "UP"})
}

func SaveUser(repo persist.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			wrapErrorAndSend(err, http.StatusBadRequest, c)
			return
		}
		var user persist.User
		err = json.Unmarshal(body, &user)
		if err != nil {
			wrapErrorAndSend(err, http.StatusInternalServerError, c)
			return
		}
		c.JSON(http.StatusOK, repo.Save(&user))
	}
}

func FindUser(repo persist.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			wrapErrorAndSend(err, http.StatusBadRequest, c)
			return
		}
		c.JSON(http.StatusOK, repo.Find(id))
	}
}
