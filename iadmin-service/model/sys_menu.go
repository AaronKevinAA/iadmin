package model

import (
	"ginProject/global"
	"sort"
)

type SysMenu struct {
	global.GvaModel
	MenuLevel  uint   `json:"menuLevel"`
	ParentId   uint   `json:"parentId" gorm:"comment:父菜单ID"` // 父菜单ID
	Path       string `json:"path" gorm:"comment:菜单对应路径"`    // 路由path
	Title      string `json:"title" gorm:"comment:菜单名"`
	Icon       string `json:"icon" gorm:"comment:菜单图标"`
	Name       string `json:"name" `
	Order      int    `json:"order"`
	ShowInMenu *bool  `json:"show_in_menu" gorm:"column:show_in_menu;type:boolean"`
	KeepAlive  *bool  `json:"keep_alive" gorm:"column:keep_alive;type:boolean"`
	Component  string `json:"Component"`
	FullPage   *bool  `json:"full_page" gorm:"column:full_page;type:boolean"`
}

type MenuList []SysMenu

type MenuItem struct {
	SysMenu
	Children []MenuItem `json:"children"`
}

type MenuTree []MenuItem

func (m MenuTree) Len() int { return len(m) }

func (m MenuTree) Swap(i, j int) { m[i], m[j] = m[j], m[i] }

func (m MenuTree) Less(i, j int) bool {
	if m[i].SysMenu.Order < m[j].SysMenu.Order {
		return true
	}
	return false
}

func (s MenuList) GetSysMenuTree(pid uint, menuLevel uint) []MenuItem {
	var menuTree []MenuItem
	// 最大层级10层
	if menuLevel == 10 {
		return menuTree
	}

	list := s.findChildren(pid, menuLevel)
	if len(list) == 0 {
		return menuTree
	}

	for _, v := range list {
		children := s.GetSysMenuTree(v.ID, menuLevel+1)
		menuTree = append(menuTree, MenuItem{v, children})
	}
	// 按照order排序
	sort.Sort(MenuTree(menuTree))
	return menuTree
}

func (s *MenuList) findChildren(pid uint, menuLevel uint) []SysMenu {
	var children []SysMenu
	for _, v := range *s {
		if v.ParentId == pid && v.MenuLevel == menuLevel {
			children = append(children, v)
		}
	}
	return children
}
