package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: noo matching record found")
var ErrDuplicate = errors.New("models: duplicate")

type AssetCategory struct {
	ID         int
	Code, Name string
}
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
