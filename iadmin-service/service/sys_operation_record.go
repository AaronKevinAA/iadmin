package service

import (
	"errors"
	"ginProject/global"
	"ginProject/model"
	"ginProject/model/request"
	"ginProject/utils"
)

//@author: [granty1](https://github.com/granty1)
//@function: CreateSysOperationRecord
//@description: 创建记录
//@param: sysOperationRecord model.SysOperationRecord
//@return: err error

func CreateSysOperationRecord(sysOperationRecord model.SysOperationRecord) (err error) {
	err = global.GvaDb.Create(&sysOperationRecord).Error
	return err
}

//@author: [granty1](https://github.com/granty1)
//@author: [piexlmax](https://github.com/piexlmax)
//@function: DeleteSysOperationRecordByIds
//@description: 批量删除记录
//@param: ids request.IdsReq
//@return: err error

func DeleteBatchSysOperationRecord(ids request.IdsReq) (err error) {
	err = global.GvaDb.Delete(&[]model.SysOperationRecord{}, "id in (?)", ids.Ids).Error
	return err
}

//@author: [granty1](https://github.com/granty1)
//@function: DeleteSysOperationRecord
//@description: 删除操作记录
//@param: sysOperationRecord model.SysOperationRecord
//@return: err error

func DeleteSysOperationRecord(sysOperationRecord model.SysOperationRecord) (err error) {
	err = global.GvaDb.Delete(&sysOperationRecord).Error
	return err
}

//@author: [granty1](https://github.com/granty1)
//@function: DeleteSysOperationRecord
//@description: 根据id获取单条操作记录
//@param: id uint
//@return: err error, sysOperationRecord model.SysOperationRecord

func GetSysOperationRecord(id uint) (err error, sysOperationRecord model.SysOperationRecord) {
	err = global.GvaDb.Where("id = ?", id).First(&sysOperationRecord).Error
	return
}

//@author: [granty1](https://github.com/granty1)
//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetSysOperationRecordInfoList
//@description: 分页获取操作记录列表
//@param: info request.SysOperationRecordSearch
//@return: err error, list interface{}, total int64

func GetSysOperationRecordInfoList(info request.SysOperationRecordSearch) (err error, list interface{}, total int64) {
	// 创建db
	db := global.GvaDb.Model(&model.SysOperationRecord{})
	var sysOperationRecords []model.SysOperationRecord
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.Method != "" {
		db = db.Where("method = ?", info.Method)
	}
	if info.Path != "" {
		db = db.Where("path LIKE ?", "%"+info.Path+"%")
	}
	if info.Status != "" {
		db = db.Where("status = ?", info.Status)
	}
	if len(info.CreatedAt) > 0 {
		db.Where("created_at between ? and ?", info.CreatedAt[0], info.CreatedAt[1])
	}
	err = db.Count(&total).Error
	err = db.Order("created_at desc").Scopes(utils.Paginate(info.Pagination.Current, info.Pagination.PageSize)).Preload("User").Find(&sysOperationRecords).Error
	return err, sysOperationRecords, total
}

func ExcelOutSysOperationRecord(outRequest request.SysOperationRecordExcelOut) (err error, excelFilePath string) {
	SysOperationRecordSearch := outRequest.SysOperationRecordSearch
	ExcelOutConfig := outRequest.ExcelOutConfig
	// 声明一个用户列表变量
	var sysOperationRecords []model.SysOperationRecord
	// 创建db
	db := global.GvaDb.Model(&model.SysOperationRecord{})
	// 如果有条件搜索 下方会自动创建搜索语句
	if SysOperationRecordSearch.Method != "" {
		db = db.Where("method = ?", SysOperationRecordSearch.Method)
	}
	if SysOperationRecordSearch.Path != "" {
		db = db.Where("path LIKE ?", "%"+SysOperationRecordSearch.Path+"%")
	}
	if SysOperationRecordSearch.Status != "" {
		db = db.Where("status = ?", SysOperationRecordSearch.Status)
	}
	if len(SysOperationRecordSearch.CreatedAt) > 0 {
		db.Where("created_at between ? and ?", SysOperationRecordSearch.CreatedAt[0], SysOperationRecordSearch.CreatedAt[1])
	}
	if !ExcelOutConfig.HasAllData {
		// 只查询当前页的数据
		err = db.Scopes(utils.Paginate(SysOperationRecordSearch.Pagination.Current, SysOperationRecordSearch.Pagination.PageSize)).Order("created_at desc").Preload("User").Find(&sysOperationRecords).Error
	} else {
		// 查询所有数据
		err = db.Order("created_at desc").Preload("User").Find(&sysOperationRecords).Error
	}
	// 如果查询出错则直接返回
	if err != nil {
		return err, ""
	}
	err, excelFilePath = utils.ExcelOut(ExcelOutConfig.HasTableHead, model.SysOperationRecordExcelOutTableHeadName(), model.SysOperationRecordExcelOutTableData(sysOperationRecords))
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
