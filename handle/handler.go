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
			wrapErrorAndSend(err, http.StatusBadRequest, c)
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
			return
		}
		c.JSON(http.StatusCreated, hideUserConfidentialFields(&user))
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
			wrapErrorAndSend(err, http.StatusBadRequest, c)
			return
		}
		user.Username = username
		pass, err := encoder.Encode(user.Password)
		if err != nil {
			wrapErrorAndSend(err, http.StatusBadRequest, c)
			return
		}
		user.Password = pass
		err = repo.Update(&user)
		if err != nil {
			wrapErrorAndSend(err, http.StatusInternalServerError, c)
			return
		}
		c.JSON(http.StatusAccepted, hideUserConfidentialFields(&user))
	}
}

func FindUser(repo persist.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, ok := ExtractUsernameContextData(c)
		if !ok {
			c.Status(http.StatusInternalServerError)
			return
		}
		user, err := repo.FindByUsername(username)
		if err != nil {
			wrapErrorAndSend(err, http.StatusInternalServerError, c)
			return
		}
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
		user, err := repo.FindByUsername(username)
		if err != nil {
			wrapErrorAndSend(err, http.StatusInternalServerError, c)
			return
		}
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
			wrapErrorAndSend(err, http.StatusBadRequest, c)
			return
		}
		post.OwnerRefer = username
		err = repo.Save(post)
		if err != nil {
			wrapErrorAndSend(err, http.StatusInternalServerError, c)
			return
		}
		c.JSON(http.StatusCreated, &post)
	}
}

func UpdatePost(repo persist.PostRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
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
		var post persist.Post
		err = json.Unmarshal(body, &post)
		if err != nil {
			wrapErrorAndSend(err, http.StatusBadRequest, c)
			return
		}
		persistPost, err := repo.Find(uint(id))
		if err != nil {
			wrapErrorAndSend(err, http.StatusInternalServerError, c)
			return
		}
		if persistPost == nil {
			wrapErrorAndSend(errors.New("no such post"), http.StatusBadRequest, c)
			return
		}
		if persistPost.OwnerRefer != username {
			wrapErrorAndSend(errors.New("post is not your"), http.StatusForbidden, c)
			return
		}
		post.ID = uint(id)
		err = repo.Update(&post)
		if err != nil {
			wrapErrorAndSend(err, http.StatusInternalServerError, c)
			return
		}
		c.JSON(http.StatusCreated, &post)
	}
}

func UpdatePostForcibly(repo persist.PostRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			c.Status(http.StatusBadRequest)
			return
		}
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			wrapErrorAndSend(err, http.StatusBadRequest, c)
			return
		}
		var post persist.Post
		err = json.Unmarshal(body, &post)
		if err != nil {
			wrapErrorAndSend(err, http.StatusBadRequest, c)
			return
		}
		persistPost, err := repo.Find(uint(id))
		if err != nil {
			wrapErrorAndSend(err, http.StatusInternalServerError, c)
			return
		}
		if persistPost == nil {
			wrapErrorAndSend(errors.New("no such post"), http.StatusBadRequest, c)
			return
		}
		post.ID = uint(id)
		err = repo.Update(&post)
		if err != nil {
			wrapErrorAndSend(err, http.StatusInternalServerError, c)
			return
		}
		c.JSON(http.StatusCreated, &post)
	}
}

func FindPost(repo persist.PostRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			c.Status(http.StatusBadRequest)
			return
		}
		post, err := repo.Find(uint(id))
		if err != nil {
			wrapErrorAndSend(err, http.StatusInternalServerError, c)
			return
		}
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
		posts, err := repo.FindAllByOwnerUsername(username)
		if err != nil {
			wrapErrorAndSend(err, http.StatusInternalServerError, c)
			return
		}
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
		posts, err := repo.FindAllByOwnerUsername(username)
		if err != nil {
			wrapErrorAndSend(err, http.StatusInternalServerError, c)
			return
		}
		c.JSON(http.StatusOK, posts)
	}
}

func DeletePost(repo persist.PostRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			c.Status(http.StatusBadRequest)
			return
		}
		username, ok := ExtractUsernameContextData(c)
		if !ok {
			c.Status(http.StatusInternalServerError)
			return
		}
		persistPost, err := repo.Find(uint(id))
		if err != nil {
			wrapErrorAndSend(err, http.StatusInternalServerError, c)
			return
		}
		if persistPost == nil {
			wrapErrorAndSend(errors.New("no such post"), http.StatusBadRequest, c)
			return
		}
		if persistPost.OwnerRefer != username {
			wrapErrorAndSend(errors.New("post is not your"), http.StatusForbidden, c)
			return
		}
		err = repo.Delete(uint(id))
		if err != nil {
			wrapErrorAndSend(err, http.StatusInternalServerError, c)
			return
		}
		c.Status(http.StatusAccepted)
	}
}

func DeletePostForcibly(repo persist.PostRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			c.Status(http.StatusBadRequest)
			return
		}
		err = repo.Delete(uint(id))
		if err != nil {
			wrapErrorAndSend(err, http.StatusInternalServerError, c)
			return
		}
		c.Status(http.StatusAccepted)
	}
}

func SaveComment(repo persist.CommentRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		postIdStr := c.Param("postId")
		postId, err := strconv.Atoi(postIdStr)
		if err != nil || postId <= 0 {
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
			wrapErrorAndSend(err, http.StatusBadRequest, c)
			return
		}
		comment.PostRefer = uint(postId)
		comment.OwnerRefer = username
		err = repo.Save(comment)
		if err != nil {
			wrapErrorAndSend(err, http.StatusInternalServerError, c)
			return
		}
		c.JSON(http.StatusOK, &comment)
	}
}

func UpdateComment(repo persist.CommentRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
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
		var comment persist.Comment
		err = json.Unmarshal(body, &comment)
		if err != nil {
			wrapErrorAndSend(err, http.StatusBadRequest, c)
			return
		}
		persistComment, err := repo.Find(uint(id))
		if err != nil {
			wrapErrorAndSend(err, http.StatusInternalServerError, c)
			return
		}
		if persistComment == nil {
			wrapErrorAndSend(errors.New("no such comment"), http.StatusBadRequest, c)
			return
		}
		if username != persistComment.OwnerRefer {
			wrapErrorAndSend(errors.New("comment is not your"), http.StatusForbidden, c)
			return
		}
		comment.ID = uint(id)
		err = repo.Update(&comment)
		if err != nil {
			wrapErrorAndSend(err, http.StatusInternalServerError, c)
			return
		}
		c.JSON(http.StatusOK, &comment)
	}
}

func UpdateCommentForcibly(repo persist.CommentRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
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
		var comment persist.Comment
		err = json.Unmarshal(body, &comment)
		if err != nil {
			wrapErrorAndSend(err, http.StatusBadRequest, c)
			return
		}
		persistComment, err := repo.Find(uint(id))
		if err != nil {
			wrapErrorAndSend(err, http.StatusInternalServerError, c)
			return
		}
		if persistComment == nil {
			wrapErrorAndSend(errors.New("no such comment"), http.StatusBadRequest, c)
			return
		}
		if username != persistComment.OwnerRefer {
			wrapErrorAndSend(errors.New("comment is not your"), http.StatusForbidden, c)
			return
		}
		comment.ID = uint(id)
		err = repo.Update(&comment)
		if err != nil {
			wrapErrorAndSend(err, http.StatusInternalServerError, c)
			return
		}
		c.JSON(http.StatusOK, &comment)
	}
}

func FindAllComments(repo persist.CommentRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, ok := ExtractUsernameContextData(c)
		if !ok {
			c.Status(http.StatusInternalServerError)
			return
		}
		comments, err := repo.FindAllByOwnerUsername(username)
		if err != nil {
			wrapErrorAndSend(err, http.StatusInternalServerError, c)
			return
		}
		c.JSON(http.StatusOK, comments)
	}
}

func FindAllCommentsByUsername(repo persist.CommentRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")
		if username == "" {
			c.Status(http.StatusBadRequest)
			return
		}
		comments, err := repo.FindAllByOwnerUsername(username)
		if err != nil {
			wrapErrorAndSend(err, http.StatusInternalServerError, c)
			return
		}
		c.JSON(http.StatusOK, comments)
	}
}

func DeleteComment(repo persist.CommentRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			c.Status(http.StatusBadRequest)
			return
		}
		username, ok := ExtractUsernameContextData(c)
		if !ok {
			c.Status(http.StatusInternalServerError)
			return
		}
		persistComment, err := repo.Find(uint(id))
		if err != nil {
			wrapErrorAndSend(err, http.StatusInternalServerError, c)
			return
		}
		if persistComment == nil {
			wrapErrorAndSend(errors.New("no such persistComment"), http.StatusBadRequest, c)
			return
		}
		if username != persistComment.OwnerRefer {
			wrapErrorAndSend(errors.New("comment is not your"), http.StatusForbidden, c)
			return
		}
		err = repo.Delete(uint(id))
		if err != nil {
			wrapErrorAndSend(err, http.StatusInternalServerError, c)
			return
		}
		c.Status(http.StatusAccepted)
	}
}

func DeleteCommentForcibly(repo persist.CommentRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			c.Status(http.StatusBadRequest)
			return
		}
		err = repo.Delete(uint(id))
		if err != nil {
			wrapErrorAndSend(err, http.StatusInternalServerError, c)
			return
		}
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
			Name string `json:"name"`
		}
		err = json.Unmarshal(body, &role)
		if err != nil {
			wrapErrorAndSend(err, http.StatusBadRequest, c)
			return
		}
		err = repo.AddRole(username, role.Name)
		if err != nil {
			wrapErrorAndSend(err, http.StatusInternalServerError, c)
			return
		}
		c.Status(http.StatusAccepted)
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
			Name string `json:"name"`
		}
		err = json.Unmarshal(body, &role)
		if err != nil {
			wrapErrorAndSend(err, http.StatusBadRequest, c)
			return
		}
		err = repo.RemoveRole(username, role.Name)
		if err != nil {
			wrapErrorAndSend(err, http.StatusInternalServerError, c)
			return
		}
		c.Status(http.StatusAccepted)
	}
}

const confidentialFieldValue = "<secret>"

func hideUserConfidentialFields(user *persist.User) *persist.User {
	user.Password = confidentialFieldValue
	return user
}
