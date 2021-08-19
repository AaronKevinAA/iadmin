package utils

import (
	"errors"
	"ginProject/global"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"time"
)

// 封装一些文件操作

// DeleteLocalFile 删除当地文件
func DeleteLocalFile(filePath string) (err error) {
	err = os.Remove(filePath)
	return err
}

// GetFileSize 获得文件大小
func GetFileSize(filePath string) (size int64) {
	file, err := os.Open(filePath)
	defer file.Close()
	if err == nil {
		fi, _ := file.Stat()
		return fi.Size()
	}
	return -1
}

// ValidFile 判断文件是否合格 大小，后缀
func ValidFile(header *multipart.FileHeader) (err error) {
	// 读取文件大小
	size := header.Size

	/* 读取配置文件中允许上传的最大文件大小 */
	maxSize := global.GvaConfig.File.MaxUploadSize
	if size > maxSize {
		return errors.New("文件过大！")
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
		return errors.New("不允许上传此类型的文件！")
	}
	return nil
}

// SaveLocalFile 保存文件到本地
// header 要保存的文件
// return saveDir 要保存到的目录
// return newFileName 重命名后的新名字
// return saveFilePath 文件保存后的保存路径
func SaveLocalFile(header *multipart.FileHeader, saveDir string) (err error, newFileName string, saveFilePath string) {
	ext := path.Ext(header.Filename)
	// 读取文件名（不带后缀）
	name := strings.TrimSuffix(header.Filename, ext)

	/* 重命名文件，防止文件重名 */
	newName := name + "_" + time.Now().Format("20060102150405") + ext
	/* 读取存放文件的路径 */
	saveFilePath = saveDir + "/" + newName
	/* 本地保存文件 */
	// 读取用户上传的文件
	in, readErr := header.Open()
	if readErr != nil {
		return errors.New("读取文件失败！"), "", ""
	}
	defer in.Close()

	// 创建预存的文件
	out, createErr := os.Create(saveFilePath)
	if createErr != nil {
		return errors.New("创建本地文件失败！"), "", ""
	}
	defer out.Close()

	// 拷贝文件
	_, copyErr := io.Copy(out, in)
	if copyErr != nil {
		return errors.New("保存本地文件失败！"), "", ""
	}
	return nil, newName, saveFilePath
}

// GetFileSuffix 读取文件后缀名
func GetFileSuffix(header *multipart.FileHeader) (suffix string) {
	ext := path.Ext(header.Filename)
	suffix = strings.ToLower(ext[1:len(ext)])
	return suffix
}
