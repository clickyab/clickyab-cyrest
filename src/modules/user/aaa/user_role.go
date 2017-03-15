package aaa

import (
	"common/assert"
	"fmt"
	"time"
)

// UserRole model
// @Model {
//		table = user_role
//		primary = false, user_id, role_id
// }
type UserRole struct {
	UserID    int64      `db:"user_id" json:"user_id"`
	RoleID    int64      `db:"role_id" json:"role_id"`
	CreatedAt *time.Time `db:"created_at" json:"created_at"`
}

// RegisterUserRole for user in transaction
func (m *Manager) RegisterUserRole(userID int64, roleIDS []int64) (userRole *UserRole, err error) {
	err = m.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			assert.Nil(m.Rollback())
		} else {
			err = m.Commit()
		}

		if err != nil {
			userRole = nil
		}
	}()

	if err != nil {
		return
	}
	err = m.DeleteUserRole(userID)
	if err != nil {
		return
	}
	roles, err := m.FindRoleByIDs(roleIDS)
	if err != nil {
		return
	}
	for i := range roles {
		UserRole := &UserRole{RoleID: roles[i].ID, UserID: userID}
		err = m.CreateUserRole(UserRole)
		if err != nil {
			break
		}
	}
	return
}

// DeleteUserRole delete
func (m *Manager) DeleteUserRole(userID int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id=?", UserRoleTableFull)
	_, err := m.GetDbMap().Exec(
		query,
		userID,
	)
	return err
}
