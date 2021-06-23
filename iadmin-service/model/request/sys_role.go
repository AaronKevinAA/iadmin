package request

import "ginProject/model"

type SysRoleMenuConfig struct {
	RoleId uint `json:"role_id"`
	Menus []model.SysMenu `json:"menus"`
}
type SysRoleApiConfig struct {
	RoleId uint `json:"role_id"`
	Apis []model.SysApi `json:"apis"`
}