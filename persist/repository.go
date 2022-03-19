package persist

type UserRepository interface {
	Save(user *User) *User
	Update(user *User) *User
	Find(id uint) *User
	FindByUsername(username string) *User
	Delete(id uint)
}

type PostRepository interface {
	Save(post *Post) *Post
	Update(post *Post) *Post
	Find(id uint) *Post
	FindAllByOwnerId(ownerId uint) []*Post
	Delete(id uint)
}

type CommentRepository interface {
	Save(comment *Comment) *Comment
	Update(comment *Comment) *Comment
	Find(id uint) *Comment
	FindAllByPostId(postId uint) []*Comment
	FindAllByOwnerId(ownerId uint) []*Comment
	Delete(id uint)
}
