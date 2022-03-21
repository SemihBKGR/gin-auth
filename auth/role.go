package auth

import "fmt"

const (
	RoleAdmin     = "ADMIN"
	RoleManager   = "MANAGER"
	RoleModerator = "MOD"
	RoleUser      = "USER"
	RoleAnonymous = "ANONYMOUS"
)

var InsertRolesQuery = fmt.Sprintf("INSERT OR IGNORE INTO roles (name) VALUES ('%s'), ('%s'), ('%s'), ('%s'), ('%s')",
	RoleAdmin, RoleManager, RoleModerator, RoleUser, RoleAnonymous)
