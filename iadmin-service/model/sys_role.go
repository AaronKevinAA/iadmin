package model

import (
	"ginProject/global"
)

type SysRole struct {
	global.GvaModel
	Name          string    `json:"name"`
	DefaultRouter string    `json:"default_router" gorm:"column:default_router;default:/index"` // 首页路径(默认/index)
	Menus         []SysMenu `json:"menus" gorm:"many2many:sys_role_menu;"`
}

func SysRoleExcelOutTableHeadName() (tableHeadName []string) {
	tableHeadName = []string{"ID", "角色名", "创建时间", "最后更新时间"}
	return tableHeadName
}

func SysRoleExcelOutTableData(roleList []SysRole) (tableData [][]string) {
	for _, role := range roleList {
		roleInfo := []string{global.Uint2String(role.ID), role.Name, global.Timestamp2DateTime(role.CreatedAt), global.Timestamp2DateTime(role.UpdatedAt)}
		tableData = append(tableData, roleInfo)
	}
	return tableData
}
