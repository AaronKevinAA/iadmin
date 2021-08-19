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

func SysMenuExcelOutTableHeadName() (tableHeadName []string) {
	tableHeadName = []string{"ID", "父菜单ID", "菜单层级", "展示标题", "路由路径", "路由名称", "前端页面路径", "展示顺序", "图标", "是否在菜单显示", "是否使用Tab缓存", "是否全屏显示", "创建时间", "最后更新时间"}
	return tableHeadName
}

func SysMenuExcelOutTableData(menuList []SysMenu) (tableData [][]string) {
	for _, menu := range menuList {
		menuInfo := []string{global.Uint2String(menu.ID), global.Uint2String(menu.ParentId), global.Uint2String(menu.MenuLevel), menu.Title, menu.Path, menu.Name, menu.Component, global.Int2String(menu.Order),
			menu.Icon, global.Bool2String(menu.ShowInMenu), global.Bool2String(menu.KeepAlive), global.Bool2String(menu.FullPage), global.Timestamp2DateTime(menu.CreatedAt), global.Timestamp2DateTime(menu.UpdatedAt)}
		tableData = append(tableData, menuInfo)
	}
	return tableData
}
