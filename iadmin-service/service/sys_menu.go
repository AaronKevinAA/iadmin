package service

import (
	"ginProject/global"
	"ginProject/model"
	"ginProject/model/request"
)

func GetSysRouteList() (err error, tree []model.MenuItem) {
	var allMenus model.MenuList
	db := global.GvaDb.Model(&model.SysMenu{})
	err = db.Order("created_at desc").Find(&allMenus).Error
	tree = allMenus.GetSysMenuTree(0, 1)
	return err, tree
}

func UpdateSysMenuInfo(req model.SysMenu) (err error, ret model.SysMenu) {
	err = global.GvaDb.Debug().Updates(&req).Error
	return err, req
}

func AddSysMenuInfo(req model.SysMenu) (err error, ret model.SysMenu) {
	err = global.GvaDb.Create(&req).Error
	return err, req
}

func DeleteBatchSysMenu(ids request.IdsReq) (err error) {
	err = global.GvaDb.Delete(&[]model.SysMenu{}, "id in (?)", ids.Ids).Error
	err = global.GvaDb.Delete(&[]model.SysMenu{}, "parent_id in (?)", ids.Ids).Error
	return err
}

func GetSysMenuByToken(roleId uint) (err error, routes model.MenuList, menus []model.MenuItem) {
	var roleMenus []model.SysRoleMenu
	err = global.GvaDb.Where("sys_role_id = ?", roleId).Preload("Menu").Find(&roleMenus).Error
	for _, v := range roleMenus {
		routes = append(routes, v.Menu)
	}
	if routes == nil {
		return err, nil, nil
	}
	menus = routes.GetSysMenuTree(0, 1)

	return err, routes, menus
}
