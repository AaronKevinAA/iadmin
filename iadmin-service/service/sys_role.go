package service

import (
	"errors"
	"ginProject/global"
	"ginProject/model"
	"ginProject/model/request"
	"ginProject/utils"
)

func GetSysRoleList(query *request.Pagination) (err error, total int64, list []model.SysRole) {
	db := global.GvaDb.Model(&model.SysRole{})
	err = db.Count(&total).Error
	if err != nil {
		return err, 0, nil
	}
	err = db.Scopes(utils.Paginate(query.Current, query.PageSize)).Order("created_at desc").Preload("Menus").Find(&list).Error
	return err, total, list
}

func UpdateSysRoleInfo(req model.SysRole) (err error, ret model.SysRole) {
	err = global.GvaDb.Updates(&req).Error
	return err, req
}

func AddSysRoleInfo(req model.SysRole) (err error, ret model.SysRole) {
	err = global.GvaDb.Create(&req).Error
	return err, req
}

func DeleteBatchSysRole(ids request.IdsReq) (err error) {
	err = global.GvaDb.Delete(&[]model.SysRole{}, "id in (?)", ids.Ids).Error
	return err
}

func UpdateSysRoleMenuConfig(config request.SysRoleMenuConfig) (err error) {
	var ret model.SysRole
	global.GvaDb.Preload("Menus").First(&ret, "id = ?", config.RoleId)
	err = global.GvaDb.Unscoped().Model(&ret).Association("Menus").Replace(&config.Menus)
	return err
}

func SetRoleDefaultRouter(data request.SysRoleDefaultRouter) (err error) {
	err = global.GvaDb.Model(model.SysRole{}).Where("id = ?", data.RoleId).Update("default_router", data.DefaultRouter).Error
	return err
}

func ExcelOutSysRole(outRequest request.SysRoleExcelOut) (err error, excelFilePath string) {
	var roleList []model.SysRole
	if !outRequest.ExcelOutConfig.HasAllData {
		// 只查询当前页的数据
		err = global.GvaDb.Scopes(utils.Paginate(outRequest.Pagination.Current, outRequest.Pagination.PageSize)).Order("created_at desc").Find(&roleList).Error
	} else {
		// 查询所有数据
		err = global.GvaDb.Order("created_at desc").Find(&roleList).Error
	}
	// 如果查询出错则直接返回
	if err != nil {
		return err, ""
	}
	err, excelFilePath = utils.ExcelOut(outRequest.ExcelOutConfig.HasTableHead, model.SysRoleExcelOutTableHeadName(), model.SysRoleExcelOutTableData(roleList))
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
