package persist

import (
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const dbName = "test.db"

var log = logrus.New()
var db *gorm.DB

func InitDatabase(callback func(db *gorm.DB)) *gorm.DB {
	if db != nil {
		if callback != nil {
			callback(db)
		}
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
	if callback != nil {
		callback(db)
	}
	return db
}

type UserSqliteRepository struct {
	db *gorm.DB
}

func NewUserSqliteRepository() *UserSqliteRepository {
	return &UserSqliteRepository{
		db: InitDatabase(nil),
	}
}

func (repo *UserSqliteRepository) Save(user *User) error {
	return repo.db.Create(user).Error
}

func (repo *UserSqliteRepository) Update(username string, user *User) error {
	return repo.db.Model(user).Updates(map[string]interface{}{"password": user.Password}).Error
}

func (repo *UserSqliteRepository) FindByUsername(username string) *User {
	var user User
	repo.db.Preload("roles").First(&user, "username = ?", username)
	return &user
}

func (repo *UserSqliteRepository) AddRole(username, role string) error {
	return repo.db.Exec(generateInsertUserRoleQuery(username, role)).Error
}

func (repo *UserSqliteRepository) RemoveRole(username, role string) error {
	return repo.db.Exec(generateDeleteUserRoleQuery(username, role)).Error
}

type PostSqliteRepository struct {
	db *gorm.DB
}

func (repo *PostSqliteRepository) Save(post *Post) error {
	return repo.db.Create(post).Error
}

func (repo *PostSqliteRepository) Update(post *Post) error {
	return repo.db.Model(post).Updates(post).Error
}

func (repo *PostSqliteRepository) Find(id uint) (*Post, error) {
	var post Post
	err := repo.db.First(&post, id).Error
	return &post, err
}

func (repo *PostSqliteRepository) FindAllByOwnerUsername(ownerUsername string) ([]*Post, error) {
	var posts []*Post
	err := repo.db.Find(&posts, "owner_refer = ?", ownerUsername).Error
	return posts, err
}

func (repo *PostSqliteRepository) Delete(id uint) error {
	var post Post
	return repo.db.Delete(&post, id).Error
}

func NewPostSqliteRepository() *PostSqliteRepository {
	return &PostSqliteRepository{
		db: InitDatabase(nil),
	}
}

type CommentSqliteRepository struct {
	db *gorm.DB
}

func (repo *CommentSqliteRepository) Save(comment *Comment) error {
	return repo.db.Create(comment).Error
}

func (repo *CommentSqliteRepository) Update(comment *Comment) error {
	return repo.db.Updates(comment).Error
}

func (repo *CommentSqliteRepository) Find(id uint) (*Comment, error) {
	var comment Comment
	err := repo.db.First(&comment, id).Error
	return &comment, err
}

func (repo *CommentSqliteRepository) FindAllByOwnerUsername(ownerUsername string) ([]*Comment, error) {
	var comments []*Comment
	err := repo.db.Find(comments, "owner_refer = ?", ownerUsername).Error
	return comments, err
}

func (repo *CommentSqliteRepository) Delete(id uint) error {
	var comment Comment
	return repo.db.Delete(&comment, id).Error
}

func NewCommentSqliteRepository() *CommentSqliteRepository {
	return &CommentSqliteRepository{
		db: InitDatabase(nil),
	}
}
