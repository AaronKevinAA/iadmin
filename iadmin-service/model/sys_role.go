package model

import "ginProject/global"

type SysRole struct {
	global.GvaModel
	Name          string    `json:"name"`
	DefaultRouter string    `json:"default_router" gorm:"column:default_router;default:/index"` // 首页路径(默认/index)
	Menus         []SysMenu `json:"menus" gorm:"many2many:sys_role_menu;"`
}
