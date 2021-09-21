package service

import (
	"ginProject/global"
	"ginProject/model"
	"ginProject/model/request"
	"ginProject/utils"
	"mime/multipart"
	"path"
	"strings"
)

// UploadFile 上传文件
func UploadFile(header *multipart.FileHeader) (file model.File, err error) {
	/* 检查文件是否合格 */
	err = utils.ValidFile(header)
	if err != nil {
		return model.File{}, err
	}
	// 读取文件名（不带后缀）
	ext := path.Ext(header.Filename)
	name := strings.TrimSuffix(header.Filename, ext)
	// 读取后缀名
	suffix := utils.GetFileSuffix(header)
	saveDir := global.GvaConfig.File.FileLocalDir
	/* 保存到本地 */
	errSave, newName, saveFilePath := utils.SaveLocalFile(header, saveDir)
	if errSave != nil {
		return model.File{}, errSave
	}
	file = model.File{
		OldName: name,
		NewName: newName,
		Suffix:  suffix,
		Size:    header.Size,
		Path:    saveFilePath,
		Url:     newName}
	/* 写入数据库 */
	err = global.GvaDb.Create(&model.File{OldName: name, NewName: newName, Suffix: suffix, Size: header.Size, Path: saveFilePath, Url: newName}).Error
	if err != nil {
		return file, err
	}
	return file, nil
}

// GetSysFileList 分页获得文件列表
func GetSysFileList(U *request.Pagination) (err error, total int64, list []model.File) {
	db := global.GvaDb.Model(&model.File{})
	err = db.Count(&total).Error
	if err != nil {
		return err, 0, nil
	}
	err = db.Scopes(utils.Paginate(U.Current, U.PageSize)).Order("created_at desc").Find(&list).Error
	return err, total, list
}

// DownloadFile 下载文件
func DownloadFile(fileID uint) (fileName string, filePath string, err error) {
	/* 查询数据库获得文件路径 */
	var file model.File
	err = global.GvaDb.Where("id = ?", fileID).Find(&file).Error
	if err != nil {
		return "", "", err
	}
	/* 拼接获得文件路径 */
	//filePath = global.GvaConfig.File.FileLocalDir + file.NewName
	filePath = file.Path
	fileName = file.OldName + "." + file.Suffix
	return fileName, filePath, nil
}

func DeleteBatchFile(ids request.IdsReq) (err error) {
	err = global.GvaDb.Delete(&[]model.File{}, "id in (?)", ids.Ids).Error
	return err
}
