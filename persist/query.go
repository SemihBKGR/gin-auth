package persist

import "fmt"

const insertUserRoleQuery = "INSERT OR IGNORE INTO user_role_join (user_id, role_id) " +
	"SELECT id as user_id, (SELECT id FROM roles WHERE name = '%s') AS role_id FROM USERS " +
	"WHERE username = '%s'"

const deleteUserRoleQuery = "DELETE FROM user_role_join " +
	"WHERE role_id = (SELECT id FROM roles WHERE name = '%s') " +
	"AND user_id = (SELECT id FROM users WHERE username = '%s')"

func generateInsertUserRoleQuery(username, role string) string {
	return fmt.Sprintf(insertUserRoleQuery, role, username)
}

func generateDeleteUserRoleQuery(username, role string) string {
	return fmt.Sprintf(deleteUserRoleQuery, role, username)
}
