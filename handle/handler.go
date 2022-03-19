package handle

import (
	"encoding/json"
	"gin-auth/auth"
	"gin-auth/persist"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
)

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, struct{ Status string }{Status: "UP"})
}

func SaveUser(repo persist.UserRepository, encoder auth.PasswordEncoder) gin.HandlerFunc {
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
		pass, err := encoder.Encode(user.Password)
		if err != nil {
			wrapErrorAndSend(err, http.StatusBadRequest, c)
			return
		}
		user.Password = pass
		c.JSON(http.StatusCreated, repo.Save(&user))
	}
}

func UpdateUser(repo persist.UserRepository, encoder auth.PasswordEncoder) gin.HandlerFunc {
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
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			wrapErrorAndSend(err, http.StatusBadRequest, c)
			return
		}
		user.ID = uint(id)
		pass, err := encoder.Encode(user.Password)
		if err != nil {
			wrapErrorAndSend(err, http.StatusBadRequest, c)
			return
		}
		user.Password = pass
		c.JSON(http.StatusAccepted, repo.Update(&user))
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
		c.JSON(http.StatusOK, repo.Find(uint(id)))
	}
}

func DeleteUser(repo persist.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			wrapErrorAndSend(err, http.StatusBadRequest, c)
			return
		}
		repo.Delete(uint(id))
		c.Status(http.StatusNoContent)
	}
}
