package model

import "ginProject/global"

type SysRole struct {
	global.GvaModel
	Name string `json:"name"`
	Menus []SysMenu	`json:"menus" gorm:"many2many:sys_role_menu;"`
	Apis []SysApi	`json:"apis" gorm:"many2many:sys_role_api;"`
}
