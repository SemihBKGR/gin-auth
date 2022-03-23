package persist

type UserRepository interface {
	Save(user *User) error
	Update(username string, user *User) error
	FindByUsername(username string) (*User, error)
	AddRole(username, role string) error
	RemoveRole(username, role string) error
}

type PostRepository interface {
	Save(post *Post) error
	Update(post *Post) error
	Find(id uint) *Post
	FindAllByOwnerUsername(ownerUsername string) []*Post
	Delete(id uint) error
}

type CommentRepository interface {
	Save(comment *Comment) error
	Update(comment *Comment) error
	Find(id uint) *Comment
	FindAllByOwnerUsername(ownerUsername string) error
	Delete(id uint) error
}
