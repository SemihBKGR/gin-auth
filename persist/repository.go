package persist

type UserRepository interface {
	save(user *User) (*User, error)
	update(user *User) (*User, error)
	find(id int) *User
	delete(id int) error
}

type PostRepository interface {
	save(post *Post) (*Post, error)
	update(post *Post) (*Post, error)
	find(id string) (*Post, error)
	findAllByOwnerId(ownerId int) []*Post
	delete(id string) error
}

type CommentRepository interface {
	save(comment *Comment) (*Comment, error)
	update(comment *Comment) (*Comment, error)
	find(id string) (*Comment, error)
	findAllByPostId(postId string) []*Comment
	findAllByOwnerId(ownerId int) []*Comment
	delete(id string) error
}
