package persist

type UserRepository interface {
	Save(user *User) *User
	Update(user *User) *User
	Find(id int) *User
	Delete(id int)
}

type PostRepository interface {
	Save(post *Post) *Post
	Update(post *Post) *Post
	Find(id string) *Post
	FindAllByOwnerId(ownerId int) []*Post
	Delete(id string)
}

type CommentRepository interface {
	Save(comment *Comment) *Comment
	Update(comment *Comment) *Comment
	Find(id string) *Comment
	FindAllByPostId(postId string) []*Comment
	FindAllByOwnerId(ownerId int) []*Comment
	Delete(id string)
}
