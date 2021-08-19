package service

import (
	"errors"
	"ginProject/global"
	"ginProject/model"
	"ginProject/model/request"
	"ginProject/utils"
	"gorm.io/gorm"
	"mime/multipart"
)

func GetSysApiList(query *request.SysApiListSearch) (err error, total int64, list []model.SysApi) {
	db := global.GvaDb.Model(&model.SysApi{}).Where("description like ? and apiGroup like ?", "%"+query.Description+"%", "%"+query.ApiGroup+"%")
	if query.Method != "" {
		db.Where("method = ?", query.Method)
	}
	if len(query.CreatedAt) > 0 {
		db.Where("created_at between ? and ?", query.CreatedAt[0], query.CreatedAt[1])
	}
	err = db.Count(&total).Error
	if err != nil {
		return err, 0, nil
	}
	err = db.Scopes(utils.Paginate(query.Pagination.Current, query.Pagination.PageSize)).Order("created_at desc").Find(&list).Error
	return err, total, list
}

func UpdateSysApiInfo(req model.SysApi) (err error) {
	var oldA model.SysApi
	err = global.GvaDb.Where("id = ?", req.ID).First(&oldA).Error
	if oldA.Path != req.Path || oldA.Method != req.Method {
		// mysql 默认字段是不区分大小写的 !!!!
		// 如果要查询的时候区分大小写，则在where后 要区分大小的字段前面加个 BINARY 就可以了 !!!
		// path 要区分大小写 method不用
		if !errors.Is(global.GvaDb.Where("BINARY path = ? AND method = ?", req.Path, req.Method).First(&model.SysApi{}).Error, gorm.ErrRecordNotFound) {
			return errors.New("存在相同api路径！")
		}
	}
	// 同步更新casbin_rule
	err = global.GvaDb.Updates(&req).Error
	if err != nil {
		return err
	} else {
		err = UpdateCasbinApi(oldA.Path, req.Path, oldA.Method, req.Method)
		if err != nil {
			return err
		} else {
			err = global.GvaDb.Save(&req).Error
		}
	}
	return err
}

func AddSysApiInfo(req model.SysApi) (err error, ret model.SysApi) {
	err = global.GvaDb.Create(&req).Error
	return err, req
}

func DeleteBatchSysApi(ids request.IdsReq) (err error) {
	err = global.GvaDb.Delete(&[]model.SysApi{}, "id in (?)", ids.Ids).Error
	return err
}

func GetSysApiTree() (err error, treeMap map[string][]model.SysApi) {
	var allApis []model.SysApi
	treeMap = make(map[string][]model.SysApi)
	// 查询所有接口
	err = global.GvaDb.Model(&model.SysApi{}).Find(&allApis).Error
	if err != nil {
		return err, nil
	}
	for _, item := range allApis {
		treeMap[item.ApiGroup] = append(treeMap[item.ApiGroup], item)
	}
	return err, treeMap
}

func ExcelOutSysApi(outRequest request.SysApiExcelOut) (err error, excelFilePath string) {
	SysApiListSearch := outRequest.SysApiListSearch
	ExcelOutConfig := outRequest.ExcelOutConfig
	// 声明一个用户列表变量
	var apiList []model.SysApi
	db := global.GvaDb.Model(&model.SysApi{}).Where("description like ? and apiGroup like ?", "%"+SysApiListSearch.Description+"%", "%"+SysApiListSearch.ApiGroup+"%")
	if SysApiListSearch.Method != "" {
		db.Where("method = ?", SysApiListSearch.Method)
	}
	if len(SysApiListSearch.CreatedAt) > 0 {
		db.Where("created_at between ? and ?", SysApiListSearch.CreatedAt[0], SysApiListSearch.CreatedAt[1])
	}
	if !ExcelOutConfig.HasAllData {
		// 只查询当前页的数据
		err = db.Scopes(utils.Paginate(SysApiListSearch.Pagination.Current, SysApiListSearch.Pagination.PageSize)).Order("created_at desc").Find(&apiList).Error
	} else {
		// 查询所有数据
		err = db.Order("created_at desc").Find(&apiList).Error
	}
	// 如果查询出错则直接返回
	if err != nil {
		return err, ""
	}
	err, excelFilePath = utils.ExcelOut(ExcelOutConfig.HasTableHead, model.SysApiExcelOutTableHeadName(), model.SysApiExcelOutTableData(apiList))
	if err != nil {
		return err, ""
	}
	// 获得文件大小
	// 3691211 305KB
	fileSize := utils.GetFileSize(excelFilePath)
	// 如果文件大于设置的允许最大下载文件大小，则返回错误
	if fileSize > global.GvaConfig.File.MaxDownloadSize {
		return errors.New("文件过大，不支持本地下载！"), excelFilePath
	}
	return nil, excelFilePath
}

func ExcelInPreviewSysApi(header *multipart.FileHeader) (err error, newFileName string, dataList [][]string, allDataCorrect bool) {
	allDataCorrect = true
	/* 检查文件是否合格 */
	err = utils.ValidFile(header)
	if err != nil {
		return err, "", nil, false
	}
	/* 保存文件到本地 */
	saveDir := global.GvaConfig.Excel.ExcelStoreDir
	errSave, newName, saveFilePath := utils.SaveLocalFile(header, saveDir)
	if errSave != nil {
		return errSave, "", nil, false
	}
	/* 读取数据列表 */
	err, dataList = utils.GetExcelDataList(saveFilePath, true)
	if err != nil {
		return err, newName, nil, false
	}

	// 对数据进行验证
	for rowIndex, _ := range dataList {
		if rowIndex == 0 {
			dataList[rowIndex] = append(dataList[rowIndex], "操作")
		} else {
			/* 检查请求方法 POST GET PUT DELETE */
			if dataList[rowIndex][2] != "POST" && dataList[rowIndex][2] != "GET" && dataList[rowIndex][2] != "PUT" && dataList[rowIndex][2] != "DELETE" {
				dataList[rowIndex][1] = "<span style='color:red'>" + dataList[rowIndex][2] + "</span>"
				dataList[rowIndex] = append(dataList[rowIndex], utils.ExcelInPreviewDataWrong)
				allDataCorrect = false
				continue
			}
			// 判断是新增还是更新
			if errOnly := global.GvaDb.Where("BINARY path = ? AND method = ?", dataList[rowIndex][1], dataList[rowIndex][2]).First(&model.SysApi{}).Error; !errors.Is(errOnly, gorm.ErrRecordNotFound) {
				/* 添加更新的标签 */
				dataList[rowIndex] = append(dataList[rowIndex], utils.ExcelInPreviewDataUpdate)
			} else {
				/* 添加新增的标签 */
				dataList[rowIndex] = append(dataList[rowIndex], utils.ExcelInPreviewDataCreate)
			}
		}
	}
	return err, newName, dataList, allDataCorrect
}

func ExcelInSysApi(request request.ExcelInRequest) (err error) {
	// 获得文件名称
	fileName := request.SaveFileName
	filePath := global.GvaConfig.Excel.ExcelStoreDir + "/" + fileName

	// excel读取文件
	errRead, dataList := utils.GetExcelDataList(filePath, false)
	if errRead != nil {
		return errRead
	}
	if dataList == nil {
		return nil
	}

	// 开启写入数据库的事务
	tx := global.GvaDb.Begin()
	for _, data := range dataList {
		// 检查手机号是否重复
		if errOnly := tx.Where("BINARY path = ? AND method = ?", data[1], data[2]).First(&model.SysApi{}).Error; !errors.Is(errOnly, gorm.ErrRecordNotFound) {
			// 执行更新操作
			err = tx.Where("BINARY path = ? AND method = ?", data[1], data[2]).Updates(&model.SysApi{Description: data[0], Path: data[1], Method: data[2], ApiGroup: data[3]}).Error
			if err != nil {
				// 更新操作遇到问题 执行回溯操作
				tx.Rollback()
				return err
			}
		} else {
			// 执行新增操作
			err = tx.Create(&model.SysApi{Description: data[0], Path: data[1], Method: data[2], ApiGroup: data[3]}).Error
			if err != nil {
				// 新增操作遇到问题 执行回溯操作
				tx.Rollback()
				return err
			}
		}
	}
	// 没有错误就提交事务
	tx.Commit()
	return nil
}
