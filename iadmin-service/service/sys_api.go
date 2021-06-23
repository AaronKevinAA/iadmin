package service

import (
	"ginProject/global"
	"ginProject/model"
	"ginProject/model/request"
	"ginProject/utils"
)

func GetSysApiList(query *request.SysApiListSearch) (err error,total int64,list []model.SysApi){
	db := global.GvaDb.Model(&model.SysApi{}).Where("description like ? and apiGroup like ?","%"+query.Description+"%","%"+query.ApiGroup+"%")
	if query.Method!=""{
		db.Where("method = ?",query.Method)
	}
	if len(query.CreatedAt)>0{
		db.Where("created_at between ? and ?",query.CreatedAt[0],query.CreatedAt[1])
	}
	err = db.Debug().Count(&total).Error
	if err != nil{
		return err,0,nil
	}
	err = db.Scopes(utils.Paginate(query.Pagination.Current,query.Pagination.PageSize)).Order("created_at desc").Find(&list).Error
	return err,total,list
}

func UpdateSysApiInfo(req model.SysApi) (err error, ret model.SysApi) {
	err = global.GvaDb.Updates(&req).Error
	return err, req
}

func AddSysApiInfo(req model.SysApi) (err error,ret model.SysApi) {
	err = global.GvaDb.Create(&req).Error
	return err, req
}

func DeleteBatchSysApi(ids request.IdsReq)(err error){
	err = global.GvaDb.Delete(&[]model.SysApi{}, "id in (?)", ids.Ids).Error
	return err
}

func GetSysApiTree()(err error,treeMap map[string][]model.SysApi){
	var allApis []model.SysApi
	treeMap = make(map[string][]model.SysApi)
	// 查询所有接口
	err = global.GvaDb.Model(&model.SysApi{}).Find(&allApis).Error
	if err!=nil{
		return err,nil
	}
	for _,item := range allApis{
		treeMap[item.ApiGroup] = append(treeMap[item.ApiGroup], item)
	}
	return err,treeMap
}