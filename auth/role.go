package auth

import "fmt"

const (
	RoleAdmin     = "ADMIN"
	RoleManager   = "MANAGER"
	RoleModerator = "MOD"
	RoleUser      = "USER"
	RoleAnonymous = "ANONYMOUS"
)

const (
	AdminUsername    = "admin"
	AdminPassword    = "password"
	AdminInsertQuery = "INSERT OR IGNORE INTO users (username,password) VALUES ('%s','%s')"
)

var InsertRolesQuery = fmt.Sprintf("INSERT OR IGNORE INTO roles (name) VALUES "+
	"('%s'), ('%s'), ('%s'), ('%s'), ('%s')",
	RoleAdmin, RoleManager, RoleModerator, RoleUser, RoleAnonymous)

var InsertAdminRoleQuery = fmt.Sprintf("INSERT OR IGNORE INTO user_role_join (user_id, role_id) "+
	"SELECT id as user_id, (SELECT id FROM roles WHERE name = '%s') AS role_id FROM USERS "+
	"WHERE username = '%s'", RoleAdmin, AdminUsername)

func GenerateInsertAdminQuery(encoder PasswordEncoder) string {
	encodedPassword, err := encoder.Encode(AdminPassword)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf(AdminInsertQuery, AdminUsername, encodedPassword)
}
