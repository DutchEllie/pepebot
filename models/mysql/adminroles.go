package mysql

import (
	"database/sql"
	"errors"

	"quenten.nl/pepebot/models"
)


type AdminRolesModel struct {
	DB *sql.DB
}

func (m *AdminRolesModel) GetAdmins() ([]*models.AdminRoles, error) {
	stmt := `SELECT id, rolename, roleid, guildid FROM adminroles`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tmp := []*models.AdminRoles{}

	for rows.Next() {
		t := &models.AdminRoles{}
		err = rows.Scan(&t.ID, &t.RoleName, &t.RoleID, &t.GuildID)
		if err != nil {
			return nil, err
		}

		tmp = append(tmp, t)
	}

	if err = rows.Err(); err != nil{
		return nil, err
	}

	return tmp, nil
}

func (m *AdminRolesModel) GetAdminRoleIDs() ([]string, error) {
	admins, err := m.GetAdmins()
	if err != nil {
		return nil, err
	}
	
	roleIDs := make([]string, 0)
	for i := 0; i < len(admins); i++ {
		roleIDs = append(roleIDs, admins[i].RoleID)
	}

	return roleIDs, nil
}

func (m *AdminRolesModel) AddAdminRole(roleName string, roleID string, guildID string) (int, error) {
	stmt := `INSERT INTO adminroles (rolename, roleid, guildid) VALUES (?, ?, ?)`

	result, err := m.DB.Exec(stmt, roleName, roleID, guildID)
	if err != nil {
		return 0, err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *AdminRolesModel) RemoveAdminRole(roleName string, roleID string, guildID string) (error) {
	stmt := `SELECT id FROM adminroles WHERE rolename = ? AND roleid = ? AND guildid = ?`

	row := m.DB.QueryRow(stmt, roleName, roleID, guildID)

	var id int
	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
			return models.ErrNoRecord
		}else{
			return err
		}
	}

	stmt = `DELETE FROM adminroles WHERE id = ?`
	result, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	if r, _ := result.RowsAffected(); r == 0 || r > 1 {
		return errors.New("either zero or more than one rows were affected")
	}

	return nil
}