package mysqlite

import (
	"ams/pkg/models"
	"ams/pkg/util"
	"database/sql"
	"fmt"
	"strings"
)

type ComputerModel struct{ DB *sql.DB }

func (m *ComputerModel) Initialize(c []*models.ComputerConfig) error {
	// skip exist ips
	rows, err := m.DB.Query("SELECT ip FROM init_computer_cs")
	if err != nil {
		return err
	}
	defer rows.Close()

	ips := []*string{}
	for rows.Next() {
		var ip string
		if err := rows.Scan(&ip); err != nil {
			return err
		}
		ips = append(ips, &ip)
	}
	if err = rows.Err(); err != nil {
		return err
	}

	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}

	// add configuration
	cs := `INSERT INTO init_computer_cs(ip,name,user_name) VALUES(?,?,?);`
	os := `INSERT INTO init_computer_os(ip,caption,version,install_date) VALUES(?,?,?,?);`
	cpu := `INSERT INTO init_computer_cpu(ip,name,number_of_cores) VALUES(?,?,?);`
	disk := `INSERT INTO init_computer_disk(ip,model,size,sn) VALUES(?,?,?,?);`
	mem := `INSERT INTO init_computer_mem(ip,manufacturer,capacity)VALUES(?,?,?);`
	net := `INSERT INTO init_computer_net(ip,description,mac)VALUES(?,?,?);`

	for _, config := range c {
		ip := config.IP
		if !util.Exist(ip, ips) {
			for _, c := range config.CS {
				_, err = tx.Exec(cs, ip, c.HostName, c.UserName)
				if err != nil {
					tx.Rollback()
					return err
				}
			}

			for _, c := range config.OS {
				var installDate = c.InstallDate // 20140319102030.000000+480
				fDate := util.DateTime(installDate)
				_, err = tx.Exec(os, ip, c.Caption, c.Version, fDate)
				if err != nil {
					tx.Rollback()
					return err
				}
			}

			for _, c := range config.CPU {
				_, err = tx.Exec(cpu, ip, c.Name, c.Cores)
				if err != nil {
					tx.Rollback()
					return err
				}
			}

			for _, c := range config.Disk {
				// convert int64 to float64
				size := float64(c.Size) / 1024 / 1024 / 1024
				s := fmt.Sprintf("%.2f", size) // 小数点后2位
				_, err = tx.Exec(disk, ip, c.Model, s, strings.Trim(c.SN, " "))
				if err != nil {
					tx.Rollback()
					return err
				}
			}

			for _, c := range config.Mem {
				// convert int64 to float64
				cap := float64(c.Capacity) / 1024 / 1024 / 1024
				_, err = tx.Exec(mem, ip, c.Manufacturer, cap)
				if err != nil {
					tx.Rollback()
					return err
				}
			}

			for _, c := range config.Net {
				_, err = tx.Exec(net, ip, c.Description, c.MAC)
				if err != nil {
					tx.Rollback()
					return err
				}
			}
		}
	}

	return tx.Commit()
}
