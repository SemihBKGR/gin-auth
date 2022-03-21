package persist

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string    `json:"username,omitempty" gorm:"unique;not null"`
	Password string    `json:"password" gorm:"size:256;not null"`
	Roles    []Role    `json:"roles,omitempty" gorm:"many2many:user_role_join"`
	Posts    []Post    `json:"posts,omitempty" gorm:"foreignKey:OwnerRefer;references:Username"`
	Comments []Comment `json:"comments,omitempty" gorm:"foreignKey:OwnerRefer;references:Username"`
}

type Role struct {
	gorm.Model
	Name string `json:"name,omitempty" gorm:"unique;not null"`
}

type Post struct {
	gorm.Model
	Content    string    `json:"content,omitempty" gorm:"non null"`
	OwnerRefer string    `json:"owner_refer,omitempty"`
	Comments   []Comment `json:"comments,omitempty" gorm:"foreignKey:PostRefer"`
}

type Comment struct {
	gorm.Model
	Content    string `json:"content,omitempty" gorm:"not null"`
	OwnerRefer string `json:"owner_id,omitempty"`
	PostRefer  uint   `json:"post_id,omitempty"`
}
