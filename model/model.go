package model

type TimeAuditable struct {
	CreateTime int64 `json:"create_time,omitempty"`
	UpdateTime int64 `json:"update_time,omitempty"`
}

type UserAuditable struct {
	CreatedBy interface{} `json:"created_by,omitempty"`
	UpdateBy  interface{} `json:"update_by,omitempty"`
}

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Roles    []Role `json:"roles,omitempty"`
	TimeAuditable
}

type Role struct {
	Id   int    `json:"id"`
	Name string `json:"name,omitempty"`
}

type Post struct {
	Id      string `json:"id"`
	Content string `json:"content,omitempty"`
	OwnerId int    `json:"owner_id,omitempty"`
	TimeAuditable
	UserAuditable
}

type Comment struct {
	Id      string `json:"id"`
	Content string `json:"content,omitempty"`
	OwnerId int    `json:"owner_id,omitempty"`
	PostId  string `json:"post_id,omitempty"`
	TimeAuditable
	UserAuditable
}
