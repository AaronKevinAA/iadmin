package model

import "ginProject/global"

type SysUser struct {
	global.GvaModel
	Phone string `json:"phone"`
	Password string `json:"-"`
	RealName string `json:"real_name"`
	Role SysRole `json:"role" gorm:"foreignKey:RoleId;AssociationForeignKey:ID;comment:用户角色"`
	RoleId uint `json:"role_id"`
}