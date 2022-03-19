package persist

import "gorm.io/gorm"

type TimeAuditable struct {
	CreateTime int64 `json:"create_time,omitempty"`
	UpdateTime int64 `json:"update_time,omitempty"`
}

type UserAuditable struct {
	CreatedBy interface{} `json:"created_by,omitempty"`
	UpdateBy  interface{} `json:"update_by,omitempty"`
}

type User struct {
	gorm.Model
	TimeAuditable
	Id       int    `json:"id"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Roles    []Role `json:"roles,omitempty" gorm:"many2many:user_role_join;"`
}

type Role struct {
	gorm.Model
	Id   int    `json:"id"`
	Name string `json:"name,omitempty"`
}

type Post struct {
	gorm.Model
	TimeAuditable
	UserAuditable
	Id      string `json:"id"`
	Content string `json:"content,omitempty"`
	Owner   User   `json:"owner_id,omitempty" gorm:"foreignKey:owner_id"`
}

type Comment struct {
	TimeAuditable
	UserAuditable
	Id      string `json:"id"`
	Content string `json:"content,omitempty"`
	Owner   User   `json:"owner_id,omitempty" gorm:"foreignKey:owner_id"`
	Post    Post   `json:"post_id,omitempty" gorm:"foreignKey:post_id"`
}
