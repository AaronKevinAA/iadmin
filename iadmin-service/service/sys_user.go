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

func Login(L *request.Login) (err error, userInter *model.SysUser) {
	/* 从redis中读取私钥 */
	privateKeyStr := global.GvaRedis.Get(L.RedisKey).Val()
	privateKey := []byte(privateKeyStr)
	/* 解密密码 */
	/* 对前端的密码做处理 这一步非常重要!!!!! */
	decryptPassword := utils.RSADecrypt(L.Password, privateKey)
	var user model.SysUser
	user.Password = utils.MD5V([]byte(decryptPassword))
	user.Phone = L.Phone
	err = global.GvaDb.Where("phone = ? AND password = ?", user.Phone, user.Password).Preload("Role").First(&user).Error
	return err, &user
}

func Register(R request.Register) (err error, userInter model.SysUser) {
	/* 从redis中读取私钥 */
	privateKeyStr := global.GvaRedis.Get(R.RedisKey).Val()
	privateKey := []byte(privateKeyStr)
	/* 解密密码 */
	/* 对前端的密码做处理 这一步非常重要!!!!! */
	decryptPassword := utils.RSADecrypt(R.Password, privateKey)
	userInter.Phone = R.Phone
	userInter.RealName = R.RealName
	userInter.Password = utils.MD5V([]byte(decryptPassword))
	// 这里有bug，注册的用户的角色不确定
	userInter.RoleId = 1
	if err = global.GvaDb.Where("phone = ?", R.Phone).First(&model.SysUser{}).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("该手机号已注册！"), userInter
	}
	err = global.GvaDb.Create(&userInter).Error
	return err, userInter
}

func FindSysUserById(id uint) (err error, user *model.SysUser) {
	var u model.SysUser
	if err = global.GvaDb.Where("id = ?", id).First(&u).Error; err != nil {
		return errors.New("用户不存在!"), &u
	}
	return nil, &u
}

func GetSysUserList(query *request.SysUserListSearch) (err error, total int64, list []model.SysUser) {
	db := global.GvaDb.Model(&model.SysUser{}).Where("phone like ? and real_name like ?", "%"+query.Phone+"%", "%"+query.RealName+"%")
	if len(query.CreatedAt) > 0 {
		db.Where("created_at between ? and ?", query.CreatedAt[0], query.CreatedAt[1])
	}
	err = db.Count(&total).Error
	if err != nil {
		return err, 0, nil
	}
	err = db.Scopes(utils.Paginate(query.Pagination.Current, query.Pagination.PageSize)).Order("created_at desc").Preload("Role").Find(&list).Error
	return err, total, list
}

func UpdateSysUserInfo(reqUser model.SysUser) (err error, user model.SysUser) {
	err = global.GvaDb.Updates(&reqUser).Error
	return err, reqUser
}

func AddSysUserInfo(reqUser model.SysUser) (err error, userInter model.SysUser) {
	var user model.SysUser
	if err = global.GvaDb.Where("phone = ?", reqUser.Phone).First(&user).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("该手机号已注册！"), userInter
	}
	reqUser.Password = utils.MD5V([]byte(utils.DefaultPassword))
	err = global.GvaDb.Create(&reqUser).Error
	return err, reqUser
}

func DeleteBatchSysUser(ids request.IdsReq) (err error) {
	err = global.GvaDb.Delete(&[]model.SysUser{}, "id in (?)", ids.Ids).Error
	return err
}

func ExcelOutSysUser(outRequest request.SysUserExcelOut) (err error, excelFilePath string) {
	SysUserListSearch := outRequest.SysUserListSearch
	ExcelOutConfig := outRequest.ExcelOutConfig
	// 声明一个用户列表变量
	var userList []model.SysUser
	db := global.GvaDb.Model(&model.SysUser{}).Where("phone like ? and real_name like ?", "%"+SysUserListSearch.Phone+"%", "%"+SysUserListSearch.RealName+"%")
	if len(SysUserListSearch.CreatedAt) > 0 {
		db.Where("created_at between ? and ?", SysUserListSearch.CreatedAt[0], SysUserListSearch.CreatedAt[1])
	}
	if !ExcelOutConfig.HasAllData {
		// 只查询当前页的数据
		err = db.Scopes(utils.Paginate(SysUserListSearch.Pagination.Current, SysUserListSearch.Pagination.PageSize)).Order("created_at desc").Preload("Role").Find(&userList).Error
	} else {
		// 查询所有数据
		err = db.Order("created_at desc").Preload("Role").Find(&userList).Error
	}
	// 如果查询出错则直接返回
	if err != nil {
		return err, ""
	}
	err, excelFilePath = utils.ExcelOut(ExcelOutConfig.HasTableHead, model.SysUserExcelOutTableHeadName(), model.SysUserExcelOutTableData(userList))
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

// ExcelInPreviewSysUser 批量导入预览
func ExcelInPreviewSysUser(header *multipart.FileHeader) (err error, newFileName string, dataList [][]string, allDataCorrect bool) {
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

	// 查询角色名称列表 便于等下检验
	var roleList []model.SysRole
	global.GvaDb.Select("name").Find(&roleList)
	var roleNameList []string
	for _, role := range roleList {
		roleNameList = append(roleNameList, role.Name)
	}

	// 对数据进行验证
	for rowIndex, _ := range dataList {
		valid := true
		if rowIndex == 0 {
			dataList[rowIndex] = append(dataList[rowIndex], "操作")
		} else {
			/* 检查手机号 */
			validPhone := utils.ValidPhone(dataList[rowIndex][1])
			if !validPhone {
				// 高亮
				dataList[rowIndex][1] = "<span style='color:red'>" + dataList[rowIndex][1] + "</span>"
				valid = false
			}
			// 检查用户角色
			if validRole := utils.IsContain(roleNameList, dataList[rowIndex][3]); !validRole {
				dataList[rowIndex][3] = "<span style='color:red'>" + dataList[rowIndex][3] + "</span>"
				valid = false
			}
			if !valid {
				allDataCorrect = false
				/* 添加错误数据的标签 */
				dataList[rowIndex] = append(dataList[rowIndex], utils.ExcelInPreviewDataWrong)
			} else {
				// 判断是新增还是更新
				if errOnly := global.GvaDb.Where("phone = ?", dataList[rowIndex][1]).First(&model.SysUser{}).Error; !errors.Is(errOnly, gorm.ErrRecordNotFound) {
					/* 添加更新的标签 */
					dataList[rowIndex] = append(dataList[rowIndex], utils.ExcelInPreviewDataUpdate)
				} else {
					/* 添加新增的标签 */
					dataList[rowIndex] = append(dataList[rowIndex], utils.ExcelInPreviewDataCreate)
				}
			}
		}
	}
	return err, newName, dataList, allDataCorrect
}

func ExcelInSysUser(request request.ExcelInRequest) (err error) {
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
	// 对数据进行加工
	for _, data := range dataList {
		// 密码加密或者设置为默认密码
		if data[2] == "" {
			data[2] = utils.MD5V([]byte(utils.DefaultPassword))
		} else {
			data[2] = utils.MD5V([]byte(data[3]))
		}
		// 将用户角色转为ID
		var role model.SysRole
		err = global.GvaDb.Select("id").Where("name = ?", data[3]).First(&role).Error
		if err != nil {
			// 有一个角色错误，则直接返回错误给前端
			return err
		}
		data[3] = global.Uint2String(role.ID)
	}

	// 开启写入数据库的事务
	tx := global.GvaDb.Begin()
	for _, data := range dataList {
		// 检查手机号是否重复
		if errOnly := tx.Where("phone = ?", data[1]).First(&model.SysUser{}).Error; !errors.Is(errOnly, gorm.ErrRecordNotFound) {
			// 执行更新操作
			err = tx.Where("phone = ?", data[1]).Updates(&model.SysUser{RealName: data[0], Password: data[2], RoleId: global.String2uint(data[3])}).Error
			if err != nil {
				// 更新操作遇到问题 执行回溯操作
				tx.Rollback()
				return err
			}
		} else {
			// 执行新增操作
			err = tx.Create(&model.SysUser{RoleId: global.String2uint(data[3]), RealName: data[0], Password: data[2], Phone: data[1]}).Error
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

func ResetSysUserPassword(userId uint) (err error) {
	err = global.GvaDb.Model(&model.SysUser{}).Where("id = ?", userId).Update("password", utils.MD5V([]byte(utils.DefaultPassword))).Error
	return err
}

func UpdatePasswordByToken(req request.UpdatePasswordByToken, userId uint) (err error) {
	/* 从redis中读取私钥 */
	privateKeyStr := global.GvaRedis.Get(req.RedisKey).Val()
	privateKey := []byte(privateKeyStr)
	/* 解密密码 */
	decryptOldPwd := utils.RSADecrypt(req.OldPassword, privateKey)
	decryptNewPwd := utils.RSADecrypt(req.NewPassword, privateKey)
	decryptConfirmPwd := utils.RSADecrypt(req.ConfirmPassword, privateKey)
	/* 检查两次密码是否一致 */
	if decryptNewPwd != decryptConfirmPwd {
		return errors.New("两次输入的密码不一致！")
	}
	/* 检查旧密码是否正确 */
	if err = global.GvaDb.Where("id = ? AND password = ?", userId, utils.MD5V([]byte(decryptOldPwd))).First(&model.SysUser{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("输入的旧密码不正确！")
	}
	err = global.GvaDb.Model(&model.SysUser{}).Where("id = ?", userId).Update("password", utils.MD5V([]byte(decryptNewPwd))).Error
	return err
}

func UpdateBasicInfoByToken(req request.SysUserBasicInfo, userId uint) (err error, newUser model.SysUser) {
	var oldUser model.SysUser
	/* 检查手机号是否已经注册 */
	if err = global.GvaDb.Where("phone = ?", req.Phone).First(&oldUser).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		if oldUser.ID != userId {
			return errors.New("该手机号已注册！"), newUser
		}
	}
	/* 检查手机号是否合格 */
	if !utils.ValidPhone(req.Phone) {
		return errors.New("手机号格式不正确！"), newUser
	}
	/* 更新用户基本信息 */
	err = global.GvaDb.Where("id = ?", userId).Updates(&model.SysUser{RealName: req.RealName, Phone: req.Phone}).Error
	global.GvaDb.Where("id = ?", userId).Preload("Role").Find(&newUser)
	return err, newUser
}
