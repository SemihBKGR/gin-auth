package persist

import (
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const dbName = "test.db"

var log = logrus.New()
var db *gorm.DB

func getDatabase() *gorm.DB {
	if db != nil {
		return db
	}
	newDb, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		log.Error(err)
	}
	log.Infoln("Database created successfully")
	db = newDb
	err = db.AutoMigrate(&User{}, &Role{}, &Post{}, &Comment{})
	if err != nil {
		log.Error(err)
	}
	log.Infoln("Database schema is ready")
	return db
}

type UserSqliteRepository struct {
	db *gorm.DB
}

func NewUserSqliteRepository() *UserSqliteRepository {
	return &UserSqliteRepository{
		db: getDatabase(),
	}
}

func (repo *UserSqliteRepository) Save(user *User) *User {
	repo.db.Create(user)
	return user
}

func (repo *UserSqliteRepository) Update(user *User) *User {
	repo.db.Model(user).Updates(user)
	return user
}

func (repo *UserSqliteRepository) Find(id uint) *User {
	var user User
	repo.db.First(&user, id)
	return &user
}

func (repo *UserSqliteRepository) FindByUsername(username string) *User {
	var user User
	repo.db.First(&user, "username = ?", username)
	return &user
}

type PostSqliteRepository struct {
	db *gorm.DB
}

func (repo *PostSqliteRepository) Save(post *Post) *Post {
	repo.db.Create(post)
	return post
}

func (repo *PostSqliteRepository) Update(post *Post) *Post {
	repo.db.Model(post).Updates(post)
	return post
}

func (repo *PostSqliteRepository) Find(id uint) *Post {
	var post Post
	repo.db.First(&post, id)
	return &post
}

func (repo *PostSqliteRepository) FindAllByOwnerId(ownerId uint) []*Post {
	var posts []*Post
	repo.db.Find(posts, "ownerId = ?", ownerId)
	return posts
}

func (repo *PostSqliteRepository) Delete(id uint) {
	var post Post
	repo.db.Delete(&post, id)
}

func NewPostSqliteRepository() *PostSqliteRepository {
	return &PostSqliteRepository{
		db: getDatabase(),
	}
}

type CommentSqliteRepository struct {
	db *gorm.DB
}

func (repo *CommentSqliteRepository) Save(comment *Comment) *Comment {
	repo.db.Create(comment)
	return comment
}

func (repo *CommentSqliteRepository) Update(comment *Comment) *Comment {
	repo.db.Model(comment).Updates(comment)
	return comment
}

func (repo *CommentSqliteRepository) Find(id uint) *Comment {
	var comment Comment
	repo.db.First(&comment, id)
	return &comment
}

func (repo *CommentSqliteRepository) FindAllByPostId(postId uint) []*Comment {
	var comments []*Comment
	repo.db.Find(comments, "postId = ?", postId)
	return comments
}

func (repo *CommentSqliteRepository) FindAllByOwnerId(ownerId uint) []*Comment {
	var comments []*Comment
	repo.db.Find(comments, "ownerId = ?", ownerId)
	return comments
}

func (repo *CommentSqliteRepository) Delete(id uint) {
	var comment Comment
	repo.db.Delete(&comment, id)
}

func NewCommentSqliteRepository() *CommentSqliteRepository {
	return &CommentSqliteRepository{
		db: getDatabase(),
	}
}
