package persist

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username,omitempty" gorm:"unique;not null"`
	Password string `json:"password,omitempty" gorm:"size:256;not null"`
	Roles    []Role `json:"roles,omitempty" gorm:"many2many:user_role_join"`
}

type Role struct {
	gorm.Model
	Name string `json:"name,omitempty" gorm:"unique;not null"`
}

type Post struct {
	gorm.Model
	Content string `json:"content,omitempty" gorm:"non null"`
	Owner   User   `json:"owner_id,omitempty" gorm:"foreignKey:owner_id;unique"`
}

type Comment struct {
	gorm.Model
	Content string `json:"content,omitempty" gorm:"not null"`
	Owner   User   `json:"owner_id,omitempty" gorm:"foreignKey:owner_id"`
	Post    Post   `json:"post_id,omitempty" gorm:"foreignKey:post_id"`
}
