package model

import "ginProject/global"

type SysMenu struct {
	global.GvaModel
	MenuLevel     uint   `json:"menuLevel"`
	ParentId      uint   `json:"parentId" gorm:"comment:父菜单ID"`     // 父菜单ID
	Path          string `json:"path" gorm:"comment:菜单对应路径"`        // 路由path
	Title         string `json:"title" gorm:"comment:菜单名"`
	Icon          string `json:"icon" gorm:"comment:菜单图标"`
	Name          string `json:"name" `
	Order         int    `json:"order"`
	HasAuthority bool   `json:"has_authority" gorm:"column:has_authority"`
	ShowInMenu bool   `json:"show_in_menu" gorm:"column:show_in_menu"`
	KeepAlive bool   `json:"keep_alive" gorm:"column:keep_alive"`
	Component string   `json:"Component" `
	FullPage bool   `json:"full_page" gorm:"column:full_page"`
}

type MenuList []SysMenu

type MenuItem struct {
	SysMenu
	Children []MenuItem `json:"children"`
}

func (m MenuList) GetSysMenuTree(pid uint,menuLevel uint) []MenuItem {
	var menuTree []MenuItem
	if menuLevel == 10 {
		return menuTree
	}

	list := m.findChildren(pid,menuLevel)
	if len(list) == 0 {
		return menuTree
	}

	for _, v := range list {
		children := m.GetSysMenuTree(v.ID, menuLevel+1)
		if v.ShowInMenu{
			menuTree = append(menuTree, MenuItem{v,children})
		}
	}

	return menuTree
}

func (m *MenuList) findChildren(pid uint, menuLevel uint) []SysMenu {
	var children []SysMenu

	for _, v := range *m {
		if v.ParentId == pid && v.MenuLevel == menuLevel && v.ShowInMenu{
			children = append(children, v)
		}
	}
	return children
}