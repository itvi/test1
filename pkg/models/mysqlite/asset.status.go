package mysqlite

import (
	"ams/pkg/models"
	"database/sql"
	"fmt"
	"strings"
)

type AssetStatusModel struct{ DB *sql.DB }

func (m *AssetStatusModel) Create(a *models.AssetStatus) error {
	q := `INSERT INTO asset_status(name)VALUES(?);`
	stmt, err := m.DB.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(a.Name)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return models.ErrDuplicate
		}
	}
	return err
}

func (m *AssetStatusModel) Edit(a *models.AssetStatus) error {
	q := `UPDATE asset_status SET name=? WHERE id=?;`
	stmt, err := m.DB.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(a.Name, a.ID)
	return err
}

func (m *AssetStatusModel) Delete(id int) error {
	q := `DELETE FROM asset_status WHERE id=?;`
	stmt, err := m.DB.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	return err
}

func (m *AssetStatusModel) GetStatuses() ([]*models.AssetStatus, error) {
	q := `SELECT id,name FROM asset_status;`
	rows, err := m.DB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []*models.AssetStatus{}

	for rows.Next() {
		c := &models.AssetStatus{}
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		data = append(data, c)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

func (m *AssetStatusModel) GetStatusByName(name string) ([]*models.AssetStatus, error) {
	s := fmt.Sprintf(`SELECT id,name FROM asset_status WHERE name LIKE %s`, "'%"+name+"%'")
	rows, err := m.DB.Query(s)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []*models.AssetStatus{}
	for rows.Next() {
		d := &models.AssetStatus{}
		if err := rows.Scan(&d.ID, &d.Name); err != nil {
			return nil, err
		}
		data = append(data, d)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

func (m *AssetStatusModel) GetStatusByID(id int) (*models.AssetStatus, error) {
	q := `SELECT id,name FROM asset_status WHERE id=?`
	c := &models.AssetStatus{}

	err := m.DB.QueryRow(q, id).Scan(&c.ID, &c.Name)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return c, nil
}
