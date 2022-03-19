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
	// Migrate the schema
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

func (u UserSqliteRepository) save(user *User) (*User, error) {
	u.db.Create(user)
	return user, nil
}

func (u UserSqliteRepository) update(user *User) (*User, error) {
	u.db.Model(user).Updates(user)
	return user, nil
}

func (u UserSqliteRepository) find(id int) *User {
	var user User
	u.db.First(&user, id)
	return &user
}

func (u UserSqliteRepository) delete(id int) error {
	var user User
	u.db.Delete(&user, id)
	return nil
}

type PostSqliteRepository struct {
	db *gorm.DB
}

func (p PostSqliteRepository) save(post *Post) (*Post, error) {
	p.db.Save(post)
	return post, nil
}

func (p PostSqliteRepository) update(post *Post) (*Post, error) {
	p.db.Model(post).Updates(post)
	return post, nil
}

func (p PostSqliteRepository) find(id string) (*Post, error) {
	var post Post
	p.db.First(&post, id)
	return &post, nil
}

func (p PostSqliteRepository) findAllByOwnerId(ownerId int) []*Post {
	var posts []*Post
	p.db.Find(posts, "ownerId = ?", ownerId)
	return posts
}

func (p PostSqliteRepository) delete(id string) error {
	var post Post
	p.db.Delete(&post, id)
	return nil
}

func NewPostSqliteRepository() *PostSqliteRepository {
	return &PostSqliteRepository{
		db: getDatabase(),
	}
}

type CommentSqliteRepository struct {
	db *gorm.DB
}

func (c CommentSqliteRepository) save(comment *Comment) (*Comment, error) {
	return
}

func (c CommentSqliteRepository) update(comment *Comment) (*Comment, error) {
	//TODO implement me
	panic("implement me")
}

func (c CommentSqliteRepository) find(id string) (*Comment, error) {
	//TODO implement me
	panic("implement me")
}

func (c CommentSqliteRepository) findAllByPostId(postId string) []*Comment {
	//TODO implement me
	panic("implement me")
}

func (c CommentSqliteRepository) findAllByOwnerId(ownerId int) []*Comment {
	//TODO implement me
	panic("implement me")
}

func (c CommentSqliteRepository) delete(id string) error {
	//TODO implement me
	panic("implement me")
}

func NewCommentSqliteRepository() *CommentSqliteRepository {
	return &CommentSqliteRepository{
		db: getDatabase(),
	}
}
