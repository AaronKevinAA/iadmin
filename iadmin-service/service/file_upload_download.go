package service

import (
	"ginProject/global"
	"ginProject/model"
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

	return file, nil
	/* 写入数据库？ */
}
