package model

import "ginProject/global"

type File struct {
	global.GvaModel
	OldName string `json:"oldName"`
	NewName string `json:"newName"`
	Suffix  string `json:"suffix"`
	Size    int64  `json:"size"`
	// 存储的路径
	Path string `json:"path"`
	// 浏览的路径
	Url string `json:"url"`
}
