package service

import (
	"errors"
	"ginProject/global"
	"ginProject/model"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"time"
)

func UploadFile(header *multipart.FileHeader) (file model.File, err error) {

	// 读取文件大小
	size := header.Size

	/* 读取配置文件中允许上传的最大文件大小 */
	maxSize := global.GvaConfig.File.MaxUploadSize
	if size > maxSize {
		return model.File{}, errors.New("文件过大！")
	}

	// 读取文件后缀
	ext := path.Ext(header.Filename)
	suffix := strings.ToLower(ext[1:len(ext)])
	/* 判断后缀是否允许，防止恶意上传 */
	typeAllow := false
	allowFileTypeList := global.GvaConfig.File.AllowFileTypeList()
	for _, allowType := range allowFileTypeList {
		if suffix == allowType {
			typeAllow = true
			break
		}
	}
	if !typeAllow {
		return model.File{}, errors.New("不允许上传此类型的文件！")
	}

	// 读取文件名（不带后缀）
	name := strings.TrimSuffix(header.Filename, ext)

	/* 重命名文件，防止文件重名 */
	newName := name + "_" + time.Now().Format("20060102150405") + ext
	/* 读取存放文件的路径 */
	savePath := global.GvaConfig.File.FileLocalUrl
	saveFilePath := savePath + "/" + newName

	/* 本地保存文件 */
	// 读取用户上传的文件
	in, readErr := header.Open()
	if readErr != nil {
		return model.File{}, errors.New("读取文件失败！")
	}
	defer in.Close()

	// 创建预存的文件
	out, createErr := os.Create(saveFilePath)
	if createErr != nil {
		return model.File{}, errors.New("创建本地文件失败！")
	}
	defer out.Close()

	// 拷贝文件
	_, copyErr := io.Copy(out, in)
	if copyErr != nil {
		return model.File{}, errors.New("保存本地文件失败！")
	}
	file = model.File{
		OldName: name,
		NewName: newName,
		Suffix:  suffix,
		Size:    size,
		Path:    saveFilePath,
		Url:     newName}

	return file, nil
	/* 写入数据库？ */
}
