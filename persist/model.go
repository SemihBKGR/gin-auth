package persist

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty" gorm:"size:256"`
	Roles    []Role `json:"roles,omitempty" gorm:"many2many:user_role_join;"`
}

type Role struct {
	gorm.Model
	Name string `json:"name,omitempty"`
}

type Post struct {
	gorm.Model
	Content string `json:"content,omitempty"`
	Owner   User   `json:"owner_id,omitempty" gorm:"foreignKey:owner_id"`
}

type Comment struct {
	gorm.Model
	Content string `json:"content,omitempty"`
	Owner   User   `json:"owner_id,omitempty" gorm:"foreignKey:owner_id"`
	Post    Post   `json:"post_id,omitempty" gorm:"foreignKey:post_id"`
}
