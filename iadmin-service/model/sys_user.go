package model

import (
	"ginProject/global"
	"strconv"
)

type SysUser struct {
	global.GvaModel
	Phone    string  `json:"phone"`
	Password string  `json:"-"`
	RealName string  `json:"real_name"`
	Role     SysRole `json:"role" gorm:"foreignKey:RoleId;AssociationForeignKey:ID;comment:用户角色"`
	RoleId   uint    `json:"role_id"`
}

// SysUserExcelOutTableHeadName 获得用户的批量导出需要的信息
func SysUserExcelOutTableHeadName() (tableHeadName []string) {
	tableHeadName = []string{"ID", "真实姓名", "手机号", "角色名", "创建时间", "最后更新时间"}
	return tableHeadName
}
func SysUserExcelOutTableData(userList []SysUser) (tableData [][]string) {
	// 循环用户列表
	for _, user := range userList {
		userInfo := []string{strconv.Itoa(int(user.ID)), user.RealName, user.Phone, user.Role.Name, global.Timestamp2DateTime(user.CreatedAt), global.Timestamp2DateTime(user.UpdatedAt)}
		tableData = append(tableData, userInfo)
	}
	return tableData
}

// SysUserExcelInTableHeadName 批量导入模板的表头
func SysUserExcelInTableHeadName() (tableHeadName []string) {
	tableHeadName = []string{"真实姓名", "手机号", "密码", "角色名"}
	return tableHeadName
}
