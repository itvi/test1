package mysqlite

import (
	"ams/pkg/models"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

type AssetModel struct{ DB *sql.DB }

func (m *AssetModel) Add(d *models.Asset) error {
	stmt, err := m.DB.Prepare(`INSERT INTO asset(number,category_code,unit,supplier,model,sn,warranty,remark)VALUES(?,?,?,?,?,?,?,?);`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	warranty := d.Warranty.Format("2006-01-02")
	_, err = stmt.Exec(d.Number, d.Category.Code, d.Unit, d.Supplier, d.Model, d.SN, warranty, d.Remark)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return models.ErrDuplicate
		}
	}
	return err
}

func (m AssetModel) Edit(d *models.Asset) error {
	stmt, err := m.DB.Prepare(`UPDATE asset SET number=?,category_code=?,supplier=?,model=?,sn=?,
	warranty=?,remark=? WHERE id=?;`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	warranty := d.Warranty.Format("2006-01-02")
	_, err = stmt.Exec(d.Number, d.Category.Code, d.Supplier, d.Model, d.SN, warranty, d.Remark, d.ID)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return models.ErrDuplicate
		}
	}
	return err
}

func (m AssetModel) Del(id int) error {
	stmt, err := m.DB.Prepare(`DELETE FROM asset WHERE id=?;`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)
	return err
}

func (m AssetModel) GetAssets() ([]*models.Asset, error) {
	q := `select a.id,
			a.number,
			b.code category_code, b.name category_name,a.unit,
			a.supplier, a.model, a.sn,a.warranty,a.remark,a.created
		from asset a
		join asset_category b on a.category_code = b.code`

	rows, err := m.DB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []*models.Asset{}
	for rows.Next() {
		d := &models.Asset{}
		if err = rows.Scan(&d.ID, &d.Number, &d.Category.Code, &d.Category.Name, &d.Unit, &d.Supplier, &d.Model, &d.SN, &d.Warranty, &d.Remark, &d.Created); err != nil {
			return nil, err
		}

		data = append(data, d)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return data, nil
}

func (m AssetModel) GetAssetByID(id int) (*models.Asset, error) {
	d := &models.Asset{}
	q := `select a.id,a.number,
				 b.id category_id, b.code category_code, b.name category_name,
				 a.supplier, a.model, a.sn,a.warranty,a.remark
		  from asset a
		  join asset_category b on a.category_code = b.code where a.id=?`

	if err := m.DB.QueryRow(q, id).Scan(&d.ID, &d.Number, &d.Category.ID, &d.Category.Code,
		&d.Category.Name, &d.Supplier, &d.Model, &d.SN, &d.Warranty, &d.Remark); err != nil {
		return nil, err
	}

	return d, nil
}

func (m AssetModel) GetAssetsByNumber(number string) ([]*models.Asset, error) {
	q := fmt.Sprintf(`select a.id,a.number,b.name
		  from asset a
		  join asset_category b on a.category_code = b.code 
          where a.number like %s`, "'%"+number+"%'")

	rows, err := m.DB.Query(q)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	assets := []*models.Asset{}
	for rows.Next() {
		a := &models.Asset{}
		if err := rows.Scan(&a.ID, &a.Number, &a.Category.Name); err != nil {
			return nil, err
		}
		assets = append(assets, a)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return assets, nil
}

// Upload files
func (m AssetModel) Upload(files []*multipart.FileHeader) error {
	var myFiles []string
	for _, file := range files {
		f, err := file.Open()
		if err != nil {
			fmt.Println(err)
			return nil
		}
		defer f.Close()

		uploadedFile := "./upload/" + file.Filename
		dst, err := os.Create(uploadedFile)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		defer dst.Close()

		// copy
		if _, err := io.Copy(dst, f); err != nil {
			fmt.Println(err)
			return nil
		}

		myFiles = append(myFiles, uploadedFile)
	}

	// insert into database
	err := m.insert(myFiles)
	if err != nil {
		//TODO：delete uploaded file！
	}
	return err
}

// insert excel files into database
func (m AssetModel) insert(files []string) error {
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}

	for _, file := range files {
		// ./upload/cpu.xlsx => cpu
		fileName := strings.Split(filepath.Base(file), ".")[0]
		f, err := excelize.OpenFile(file)
		if err != nil {
			fmt.Println(err)
			return err
		}

		rows, err := f.GetRows("Sheet1")
		if err != nil {
			fmt.Println(err)
			return err
		}

		var values string

		switch strings.ToUpper(fileName) {

		case "资产":
			values = ""
			for i, row := range rows {
				if i == 0 {
					continue
				}
				values += fmt.Sprintf("(%s,%s,%s,%s,%s,%s,%s,%s),", "'"+row[0]+"'", "'"+row[1]+"'", "'"+row[2]+"'", "'"+row[3]+"'", "'"+row[4]+"'", "'"+row[5]+"'", "'"+row[6]+"'", "'"+row[7]+"'")
			}
			sql := fmt.Sprintf("INSERT INTO asset(number,category_code,unit,supplier,model,sn,warranty,remark) VALUES %s;", strings.TrimSuffix(values, ","))
			_, err = tx.Exec(sql)
			if err != nil {
				tx.Rollback()
				return err
			}

		default:
			e := errors.New("文件不正确！")
			tx.Rollback()
			return e
		}
	}
	return tx.Commit()
}

// Exist check asset exist or not
func (m AssetModel) Exist(assetNumber string) bool {
	var num int
	if err := m.DB.QueryRow(`SELECT COUNT(1) FROM asset WHERE number=?`, assetNumber).Scan(&num); err != nil {
		fmt.Println(err)
		return false
	}
	if num > 0 {
		return true
	}
	return false
}
