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
		err = repo.Save(&user)
		if err != nil {
			wrapErrorAndSend(err, http.StatusInternalServerError, c)
		} else {
			c.JSON(http.StatusCreated, hideUserConfidentialFields(&user))
		}
	}
}

func UpdateUser(repo persist.UserRepository, encoder auth.PasswordEncoder) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, ok := ExtractUsernameContextData(c)
		if !ok {
			wrapErrorAndSend(errors.New("context data does not contains username"), http.StatusInternalServerError, c)
			return
		}
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
		persistUser := repo.FindByUsername(username)
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

func FindUserByUsername(repo persist.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")
		if username != "" {
			c.Status(http.StatusBadRequest)
			return
		}
		user := repo.FindByUsername(username)
		c.JSON(http.StatusOK, hideUserConfidentialFields(user))
	}
}

func SavePost(repo persist.PostRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, ok := ExtractUsernameContextData(c)
		if !ok {
			c.Status(http.StatusInternalServerError)
			return
		}
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
		post.OwnerRefer = username
		c.JSON(http.StatusCreated, repo.Save(post))
	}
}

func UpdatePost(repo persist.PostRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		username, ok := ExtractUsernameContextData(c)
		if !ok {
			c.Status(http.StatusInternalServerError)
			return
		}
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
		persistPost := repo.Find(uint(id))
		if persistPost == nil {
			wrapErrorAndSend(errors.New("no such post"), http.StatusBadRequest, c)
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

func UpdatePostForcibly(repo persist.PostRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
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
		persistPost := repo.Find(uint(id))
		if persistPost == nil {
			wrapErrorAndSend(errors.New("no such post"), http.StatusBadRequest, c)
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
		username, ok := ExtractUsernameContextData(c)
		if !ok {
			c.Status(http.StatusInternalServerError)
			return
		}
		persistPost := repo.Find(uint(id))
		if persistPost == nil {
			wrapErrorAndSend(errors.New("no such post"), http.StatusBadRequest, c)
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

func SaveComment(repo persist.CommentRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, ok := ExtractUsernameContextData(c)
		if !ok {
			c.Status(http.StatusInternalServerError)
			return
		}
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			wrapErrorAndSend(err, http.StatusBadRequest, c)
			return
		}
		var comment *persist.Comment
		err = json.Unmarshal(body, &comment)
		if err != nil {
			wrapErrorAndSend(err, http.StatusInternalServerError, c)
			return
		}
		comment.OwnerRefer = username
		c.JSON(http.StatusOK, repo.Save(comment))
	}
}

func UpdateComment(repo persist.CommentRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		username, ok := ExtractUsernameContextData(c)
		if !ok {
			c.Status(http.StatusInternalServerError)
			return
		}
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			wrapErrorAndSend(err, http.StatusBadRequest, c)
			return
		}
		var comment *persist.Comment
		err = json.Unmarshal(body, &comment)
		if err != nil {
			wrapErrorAndSend(err, http.StatusInternalServerError, c)
			return
		}
		persistComment := repo.Find(uint(id))
		if persistComment == nil {
			wrapErrorAndSend(errors.New("no such comment"), http.StatusBadRequest, c)
			return
		}
		if username != persistComment.OwnerRefer {
			wrapErrorAndSend(errors.New("comment is not your"), http.StatusForbidden, c)
			return
		}
		persistComment.Content = comment.Content
		c.JSON(http.StatusOK, repo.Save(persistComment))
	}
}

func UpdateCommentForcibly(repo persist.CommentRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		username, ok := ExtractUsernameContextData(c)
		if !ok {
			c.Status(http.StatusInternalServerError)
			return
		}
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			wrapErrorAndSend(err, http.StatusBadRequest, c)
			return
		}
		var comment *persist.Comment
		err = json.Unmarshal(body, &comment)
		if err != nil {
			wrapErrorAndSend(err, http.StatusInternalServerError, c)
			return
		}
		persistComment := repo.Find(uint(id))
		if persistComment == nil {
			wrapErrorAndSend(errors.New("no such comment"), http.StatusBadRequest, c)
			return
		}
		if username != persistComment.OwnerRefer {
			wrapErrorAndSend(errors.New("comment is not your"), http.StatusForbidden, c)
			return
		}
		persistComment.Content = comment.Content
		c.JSON(http.StatusOK, repo.Save(persistComment))
	}
}

func FindComments(repo persist.CommentRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, ok := ExtractUsernameContextData(c)
		if !ok {
			c.Status(http.StatusInternalServerError)
			return
		}
		comments := repo.FindAllByOwnerUsername(username)
		c.JSON(http.StatusOK, comments)
	}
}

func DeleteComment(repo persist.CommentRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		username, ok := ExtractUsernameContextData(c)
		if !ok {
			c.Status(http.StatusInternalServerError)
			return
		}
		persistComment := repo.Find(uint(id))
		if persistComment == nil {
			wrapErrorAndSend(errors.New("no such persistComment"), http.StatusBadRequest, c)
			return
		}
		if username != persistComment.OwnerRefer {
			wrapErrorAndSend(errors.New("comment is not your"), http.StatusForbidden, c)
			return
		}
		repo.Delete(uint(id))
		c.Status(http.StatusAccepted)
	}
}

func DeleteCommentForcibly(repo persist.CommentRepository) gin.HandlerFunc {
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

func AddRole(repo persist.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")
		if username == "" {
			c.Status(http.StatusBadRequest)
			return
		}
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			wrapErrorAndSend(err, http.StatusBadRequest, c)
			return
		}
		var role struct {
			Name string
		}
		err = json.Unmarshal(body, &role)
		if err != nil {
			wrapErrorAndSend(err, http.StatusInternalServerError, c)
			return
		}
		repo.AddRole(username, role.Name)
	}
}

func RemoveRole(repo persist.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")
		if username == "" {
			c.Status(http.StatusBadRequest)
			return
		}
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			wrapErrorAndSend(err, http.StatusBadRequest, c)
			return
		}
		var role struct {
			Name string
		}
		err = json.Unmarshal(body, &role)
		if err != nil {
			wrapErrorAndSend(err, http.StatusInternalServerError, c)
			return
		}
		repo.RemoveRole(username, role.Name)
	}
}

const confidentialFieldValue = "<secret>"

func hideUserConfidentialFields(user *persist.User) *persist.User {
	user.Password = confidentialFieldValue
	return user
}
