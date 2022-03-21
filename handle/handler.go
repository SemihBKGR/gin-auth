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
		if !ok {
			wrapErrorAndSend(errors.New("incorrect credentials"), http.StatusUnauthorized, c)
			return
		}
		token := authTokenPrefix + jwtService.GenerateToken(user)
		c.Data(http.StatusAccepted, "text/plain", []byte(token))
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
		persistUser := repo.Save(&user)
		c.JSON(http.StatusCreated, hideUserConfidentialFields(persistUser))
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
		id, ok := ExtractIdContextData(c)
		if !ok {
			wrapErrorAndSend(errors.New("context data does not contains id"), http.StatusInternalServerError, c)
			return
		}
		persistUser := repo.Find(id)
		pass, err := encoder.Encode(user.Password)
		if err != nil {
			wrapErrorAndSend(err, http.StatusBadRequest, c)
			return
		}
		persistUser.Password = pass
		updateUser := repo.Update(persistUser)
		c.JSON(http.StatusAccepted, hideUserConfidentialFields(updateUser))
	}
}

func FindUser(repo persist.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, ok := ExtractUsernameContextData(c)
		if !ok {
			c.Status(http.StatusInternalServerError)
			return
		}
		user := repo.FindByUsername(username)
		c.JSON(http.StatusOK, hideUserConfidentialFields(user))
	}
}

func FindUserById(repo persist.UserRepository) gin.HandlerFunc {
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
		username, ok := ExtractUsernameContextData(c)
		if !ok {
			c.Status(http.StatusInternalServerError)
			return
		}
		post.OwnerRefer = username
		c.JSON(http.StatusCreated, repo.Save(post))
	}
}

func UpdatePost(repo persist.PostRepository) gin.HandlerFunc {
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
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		persistPost := repo.Find(uint(id))
		if persistPost == nil {
			wrapErrorAndSend(errors.New("no such post"), http.StatusBadRequest, c)
			return
		}
		username, ok := ExtractUsernameContextData(c)
		if !ok {
			c.Status(http.StatusInternalServerError)
			return
		}
		if persistPost.OwnerRefer != username {
			wrapErrorAndSend(errors.New("post is not your"), http.StatusForbidden, c)
			return
		}
		persistPost.Content = post.Content
		c.JSON(http.StatusCreated, repo.Save(persistPost))
	}
}

func FindPost(repo persist.PostRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		post := repo.Find(uint(id))
		c.JSON(http.StatusOK, post)
	}
}

func FindAllPosts(repo persist.PostRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, ok := ExtractUsernameContextData(c)
		if !ok {
			c.Status(http.StatusInternalServerError)
			return
		}
		posts := repo.FindAllByOwnerUsername(username)
		c.JSON(http.StatusOK, posts)
	}
}

func FindAllPostsByUsername(repo persist.PostRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")
		if username == "" {
			c.Status(http.StatusBadRequest)
			return
		}
		posts := repo.FindAllByOwnerUsername(username)
		c.JSON(http.StatusOK, posts)
	}
}

func DeletePost(repo persist.PostRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		persistPost := repo.Find(uint(id))
		if persistPost == nil {
			wrapErrorAndSend(errors.New("no such post"), http.StatusBadRequest, c)
			return
		}
		username, ok := ExtractUsernameContextData(c)
		if !ok {
			c.Status(http.StatusInternalServerError)
			return
		}
		if persistPost.OwnerRefer != username {
			wrapErrorAndSend(errors.New("post is not your"), http.StatusForbidden, c)
			return
		}
		repo.Delete(uint(id))
		c.Status(http.StatusAccepted)
	}
}

func DeletePostForcibly(repo persist.PostRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		repo.Delete(uint(id))
		c.Status(http.StatusAccepted)
	}
}

const confidentialFieldValue = "<secret>"

func hideUserConfidentialFields(user *persist.User) *persist.User {
	user.Password = confidentialFieldValue
	return user
}
