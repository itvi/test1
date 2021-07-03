package mysqlite

import (
	"ams/pkg/models"
	"database/sql"
	"fmt"
	"strings"
)

type AssetCategoryModel struct{ DB *sql.DB }

func (m *AssetCategoryModel) Create(a *models.AssetCategory) error {
	q := `INSERT INTO asset_category(code,name)VALUES(?,?);`
	stmt, err := m.DB.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(a.Code, a.Name)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return models.ErrDuplicate
		}
	}
	return err
}

func (m *AssetCategoryModel) Edit(a *models.AssetCategory) error {
	q := `UPDATE asset_category SET code=?,name=? WHERE id=?;`
	stmt, err := m.DB.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(a.Code, a.Name, a.ID)
	return err
}

func (m *AssetCategoryModel) Delete(id int) error {
	q := `DELETE FROM asset_category WHERE id=?;`
	stmt, err := m.DB.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	return err
}

func (m *AssetCategoryModel) GetCategories() ([]*models.AssetCategory, error) {
	q := `SELECT id,code,name FROM asset_category;`
	rows, err := m.DB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []*models.AssetCategory{}

	for rows.Next() {
		c := &models.AssetCategory{}
		if err := rows.Scan(&c.ID, &c.Code, &c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (m *AssetCategoryModel) GetCategoriesByName(name string) ([]*models.AssetCategory, error) {
	s := fmt.Sprintf(`SELECT id,code,name FROM asset_category WHERE name LIKE %s`, "'%"+name+"%'")
	rows, err := m.DB.Query(s)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []*models.AssetCategory{}
	for rows.Next() {
		d := &models.AssetCategory{}
		if err := rows.Scan(&d.ID, &d.Code, &d.Name); err != nil {
			return nil, err
		}
		data = append(data, d)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

func (m *AssetCategoryModel) GetCategoryByID(id int) (*models.AssetCategory, error) {
	q := `SELECT id,code,name FROM asset_category WHERE id=?`
	c := &models.AssetCategory{}

	err := m.DB.QueryRow(q, id).Scan(&c.ID, &c.Code, &c.Name)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return c, nil
}

func (m *AssetCategoryModel) GetCategoryByCode(code string) (*models.AssetCategory, error) {
	q := `SELECT id,code,name FROM asset_category WHERE code=?`
	c := &models.AssetCategory{}

	err := m.DB.QueryRow(q, code).Scan(&c.ID, &c.Code, &c.Name)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return c, nil
}
