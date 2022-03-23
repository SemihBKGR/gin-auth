package persist

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string    `json:"username" gorm:"unique;not null"`
	Password string    `json:"password" gorm:"size:256;not null"`
	Roles    []Role    `json:"roles" gorm:"many2many:user_role_join"`
	Posts    []Post    `json:"posts,omitempty" gorm:"foreignKey:OwnerRefer;references:Username"`
	Comments []Comment `json:"comments,omitempty" gorm:"foreignKey:OwnerRefer;references:Username"`
}

type Role struct {
	gorm.Model
	Name string `json:"name" gorm:"unique;not null"`
}

type Post struct {
	gorm.Model
	Content    string    `json:"content" gorm:"non null"`
	OwnerRefer string    `json:"owner_refer"`
	Comments   []Comment `json:"comments" gorm:"foreignKey:PostRefer"`
}

type Comment struct {
	gorm.Model
	Content    string `json:"content" gorm:"not null"`
	OwnerRefer string `json:"owner_id"`
	PostRefer  uint   `json:"post_id"`
}
