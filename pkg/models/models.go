package models

import (
	"ams/pkg/util"
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: noo matching record found")
var ErrDuplicate = errors.New("models: duplicate")
var ErrInvalidCredentials = errors.New("models: invalid credentials")

type Asset struct {
	ID       int
	Number   string // 设备编号 Device Number
	Unit     string // 计量单位
	Supplier string
	Model    string    // 设备型号
	SN       string    // 设备序列号
	Warranty time.Time // 质保期
	Category AssetCategory
	Remark   string
	Created  time.Time
}

type AssetCategory struct {
	ID         int
	Code, Name string
}

type AssetStatus struct {
	ID   int
	Name string
}

type User struct {
	ID             int
	SN             string
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type Department struct {
	ID   int
	Code string
	Name string
}

type AssetManagement struct {
	ID           int
	Asset        Asset
	Mvt          string
	Qty          int
	FromLoc      string
	ToLoc        string
	FromEmployee string
	ToEmployee   string
	DocumentDate time.Time
	Remark       string
}

type ComputerConfig struct {
	IP string `json:"ip"`
	CS []struct {
		HostName string `json:"hostName"`
		UserName string `json:"userName"`
	} `json:"cs"`
	OS []struct {
		Caption        string `json:"caption"`
		Version        string `json:"version"`
		InstallDate    string `json:"installDate"`
		OSArchitecture string `json:"osArchitecture"`
	} `json:"os"`
	CPU []struct {
		Name  string `json:"name"`
		Cores int    `json:"cores"`
	} `json:"cpu"`
	Disk []struct {
		Model string `json:"model"`
		Size  int64  `json:"size"`
		SN    string `json:"sn"`
	} `json:"disk"`
	Mem []struct {
		Manufacturer string `json:"manufacturer"`
		Capacity     int64  `json:"capacity"`
	} `json:"mem"`
	Net []struct {
		Description string `json:"description"`
		MAC         string `json:"mac"`
	} `json:"net"`
	Product []struct {
		Name        string `json:"name"`
		Version     string `json:"version"`
		Vendor      string `json:"vendor"`
		InstallDate string `json:"installDate"`
	} `json:"product"`
}

type AssetMovAndConfig struct {
	AssetMov struct {
		Number       string        `json:"assetNumber"`
		Mvt          string        `json:"mvt"`
		Qty          int           `json:"qty"`
		FromLoc      string        `json:"fromLoc"`
		ToLoc        string        `json:"toLoc"`
		FromEmployee string        `json:"fromEmployee"`
		ToEmployee   string        `json:"toEmployee"`
		DocumentDate util.JsonDate `json:"documentDate"`
		Remark       string        `json:"remark"`
	} `json:"assetMov"`
	Config         bool            `json:"config"`
	ComputerConfig *ComputerConfig `json:"computerConfig"`
}
