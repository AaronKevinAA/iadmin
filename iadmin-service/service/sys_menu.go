package service

import (
	"errors"
	"ginProject/global"
	"ginProject/model"
	"ginProject/model/request"
	"ginProject/utils"
)

func GetSysRouteList() (err error, tree []model.MenuItem) {
	var allMenus model.MenuList
	db := global.GvaDb.Model(&model.SysMenu{})
	err = db.Order("created_at desc").Find(&allMenus).Error
	tree = allMenus.GetSysMenuTree(0, 1)
	return err, tree
}

func UpdateSysMenuInfo(req model.SysMenu) (err error, ret model.SysMenu) {
	err = global.GvaDb.Updates(&req).Error
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

func ExcelOutSysMenu(outRequest request.ExcelOutRequest) (err error, excelFilePath string) {
	var menuList []model.SysMenu
	// 如果查询出错则直接返回
	err = global.GvaDb.Order("created_at desc").Find(&menuList).Error
	if err != nil {
		return err, ""
	}
	err, excelFilePath = utils.ExcelOut(outRequest.HasTableHead, model.SysMenuExcelOutTableHeadName(), model.SysMenuExcelOutTableData(menuList))
	if err != nil {
		return err, ""
	}
	// 获得文件大小
	fileSize := utils.GetFileSize(excelFilePath)
	// 如果文件大于设置的允许最大下载文件大小，则返回错误
	if fileSize > global.GvaConfig.File.MaxDownloadSize {
		return errors.New("文件过大，不支持本地下载！"), excelFilePath
	}
	return nil, excelFilePath
}
