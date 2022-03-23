package persist

type UserRepository interface {
	Save(user *User) error
	Update(user *User) error
	FindByUsername(username string) (*User, error)
	AddRole(username, role string) error
	RemoveRole(username, role string) error
}

type PostRepository interface {
	Save(post *Post) error
	Update(post *Post) error
	Find(id uint) (*Post, error)
	FindAllByOwnerUsername(ownerUsername string) ([]*Post, error)
	Delete(id uint) error
}

type CommentRepository interface {
	Save(comment *Comment) error
	Update(comment *Comment) error
	Find(id uint) (*Comment, error)
	FindAllByOwnerUsername(ownerUsername string) ([]*Comment, error)
	Delete(id uint) error
}
