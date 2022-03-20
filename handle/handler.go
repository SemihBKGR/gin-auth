package handle

import (
	"encoding/json"
	"errors"
	"gin-auth/auth"
	"gin-auth/auth/jwt"
	"gin-auth/persist"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
)

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, struct{ Status string }{Status: "UP"})
}

func Login(loginService auth.LoginService, jwtService jwt.JwtService) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			wrapErrorAndSend(err, http.StatusBadRequest, c)
			return
		}
		credentials := &struct {
			Username string
			Password string
		}{}
		err = json.Unmarshal(body, credentials)
		if err != nil {
			wrapErrorAndSend(err, http.StatusBadRequest, c)
			return
		}
		user, ok := loginService.Login(credentials.Username, credentials.Password)
		if user != nil {
			if ok {
				token := authTokenPrefix + jwtService.GenerateToken(user)
				c.Data(http.StatusAccepted, "text/plain", []byte(token))
			} else {
				wrapErrorAndSend(errors.New("wrong password"), http.StatusUnauthorized, c)
			}
		} else {
			wrapErrorAndSend(errors.New("no such user"), http.StatusUnauthorized, c)
		}
	}
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

func SavePost(repo persist.PostRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			wrapErrorAndSend(err, http.StatusBadRequest, c)
			return
		}
		var post *persist.Post
		err = json.Unmarshal(body, &post)
		if err != nil {
			wrapErrorAndSend(err, http.StatusInternalServerError, c)
			return
		}
		c.JSON(http.StatusCreated, repo.Save(post))
	}
}
