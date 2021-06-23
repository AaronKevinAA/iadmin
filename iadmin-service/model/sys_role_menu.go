package model

type SysRoleMenu struct {
	//SysMenu
	SysMenuId uint `json:"-"`
	SysRoleId uint	`json:"-"`
	Menu SysMenu `json:"menu" gorm:"foreignKey:SysMenuId;AssociationForeignKey:ID;comment:用户角色"`
}
