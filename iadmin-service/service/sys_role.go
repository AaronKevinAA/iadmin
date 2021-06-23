package service

import (
	"ginProject/global"
	"ginProject/model"
	"ginProject/model/request"
	"ginProject/utils"
)

func GetSysRoleList(query *request.Pagination) (err error,total int64,list []model.SysRole){
	db := global.GvaDb.Model(&model.SysRole{})
	err = db.Debug().Count(&total).Error
	if err != nil{
		return err,0,nil
	}
	err = db.Scopes(utils.Paginate(query.Current,query.PageSize)).Order("created_at desc").Preload("Menus").Preload("Apis").Find(&list).Error

	return err,total,list
}

func UpdateSysRoleInfo(req model.SysRole) (err error, ret model.SysRole) {
	err = global.GvaDb.Updates(&req).Error
	return err, req
}

func AddSysRoleInfo(req model.SysRole) (err error,ret model.SysRole) {
	err = global.GvaDb.Create(&req).Error
	return err, req
}

func DeleteBatchSysRole(ids request.IdsReq)(err error){
	err = global.GvaDb.Delete(&[]model.SysRole{}, "id in (?)", ids.Ids).Error
	return err
}

func UpdateSysRoleMenuConfig(config request.SysRoleMenuConfig)(err error){
	var ret model.SysRole
	global.GvaDb.Preload("Menus").First(&ret, "id = ?", config.RoleId)
	err = global.GvaDb.Unscoped().Model(&ret).Association("Menus").Replace(&config.Menus)
	return err
}

func UpdateSysRoleApiConfig(config request.SysRoleApiConfig)(err error){
	var ret model.SysRole
	global.GvaDb.Preload("Apis").First(&ret, "id = ?", config.RoleId)
	err = global.GvaDb.Unscoped().Model(&ret).Association("Apis").Replace(&config.Apis)
	return err
}