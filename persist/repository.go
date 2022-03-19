package persist

import (
	"gin-auth/model"
)

type UserRepository interface {
	save(user *model.User) (*model.User, error)
	update(user *model.User) (*model.User, error)
	find(id int) *model.User
	delete(id int) error
}

type RoleRepository interface {
	findAll() []*model.Role
	addRoleToUser(roleId, userId int)
}

type PostRepository interface {
	save(post *model.Post) (*model.Post, error)
	update(post *model.Post) (*model.Post, error)
	find(id string) (*model.Post, error)
	findAllByUserId(userId int) []*model.Post
	delete(id string) error
}

type CommentRepository interface {
	save(comment *model.Comment) (*model.Comment, error)
	update(comment *model.Comment) (*model.Comment, error)
	find(id string) (*model.Comment, error)
	findAllByPost(id string) []*model.Comment
	findAllByUser(id int) []*model.Comment
	delete(id string) error
}
