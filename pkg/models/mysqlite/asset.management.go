package mysqlite

import (
	"ams/pkg/models"
	"ams/pkg/util"
	"database/sql"
	"time"
)

type AssetManagementMode struct{ DB *sql.DB }

func (m *AssetManagementMode) Add(a *models.AssetManagement) error {
	stmt, err := m.DB.Prepare(`INSERT INTO asset_mov(a_number,mvt,qty,from_loc,to_loc,from_employee,to_employee,doc_date) VALUES(?,?,?,?,?,?,?,?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	docDate := a.DocumentDate.Format("2006-01-02")
	_, err = stmt.Exec(a.Asset.Number, a.Mvt, a.Qty, a.FromLoc, a.ToLoc, a.FromEmployee, a.ToEmployee, docDate)
	return err
}

// AddConfig add asset movement and computer configuration
func (m *AssetManagementMode) AddMovAndConfig(c *models.AssetMovAndConfig) error {
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}

	// device warehouse
	sqlMov := `INSERT INTO asset_mov(a_number,mvt,qty,from_loc,to_loc,from_employee,to_employee,doc_date,remark) VALUES(?,?,?,?,?,?,?,?,?)`
	docDate := time.Time(c.AssetMov.DocumentDate).Format("2006-01-02")
	_, err = tx.Exec(sqlMov, &c.AssetMov.Number, &c.AssetMov.Mvt, &c.AssetMov.Qty, &c.AssetMov.FromLoc, &c.AssetMov.ToLoc,
		&c.AssetMov.FromEmployee, &c.AssetMov.ToEmployee, docDate, &c.AssetMov.Remark)
	if err != nil {
		tx.Rollback()
		return err
	}

	if c.Config { // with configuration
		// desktop computer information
		cs := `INSERT INTO computer_cs(asset_number,name,user_name) VALUES(?,?,?);`
		os := `INSERT INTO computer_os(asset_number,caption,version,install_date) VALUES(?,?,?,?);`
		cpu := `INSERT INTO computer_cpu(asset_number,name,number_of_cores) VALUES(?,?,?);`
		disk := `INSERT INTO computer_disk(asset_number,model,size,sn) VALUES(?,?,?,?);`
		mem := `INSERT INTO computer_mem(asset_number,manufacturer,capacity)VALUES(?,?,?);`
		net := `INSERT INTO computer_net(asset_number,description,mac)VALUES(?,?,?);`

		for _, p := range c.ComputerConfig.CS {
			_, err = tx.Exec(cs, &c.AssetMov.Number, p.HostName, p.UserName)
			if err != nil {
				tx.Rollback()
				return err
			}
		}

		for _, p := range c.ComputerConfig.OS {
			// format install data (default is 20190218165055.000000+480)
			installDate := util.DateTime(p.InstallDate)
			_, err = tx.Exec(os, &c.AssetMov.Number, p.Caption, p.Version, installDate)
			if err != nil {
				tx.Rollback()
				return err
			}
		}

		for _, p := range c.ComputerConfig.CPU {
			_, err = tx.Exec(cpu, &c.AssetMov.Number, p.Name, p.Cores)
			if err != nil {
				tx.Rollback()
				return err
			}
		}

		for _, d := range c.ComputerConfig.Disk {
			_, err = tx.Exec(disk, &c.AssetMov.Number, d.Model, d.Size, d.SN)
			if err != nil {
				tx.Rollback()
				return err
			}
		}

		for _, m := range c.ComputerConfig.Mem {
			_, err = tx.Exec(mem, &c.AssetMov.Number, m.Manufacturer, m.Capacity)
			if err != nil {
				tx.Rollback()
				return err
			}
		}

		for _, n := range c.ComputerConfig.Net {
			_, err = tx.Exec(net, &c.AssetMov.Number, n.Description, n.MAC)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	return tx.Commit()
}
