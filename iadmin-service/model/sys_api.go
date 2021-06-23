package model

import "ginProject/global"

type SysApi struct {
	global.GvaModel
	Path string	`json:"path"`
	Description string `json:"description"`
	ApiGroup string	`json:"apiGroup" gorm:"column:apiGroup;"`
	Method string `json:"method"`
}
