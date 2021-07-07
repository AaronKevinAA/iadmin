package config

import "strings"

type File struct {
	MaxUploadSize int64  `mapstructure:"max-upload-size" json:"maxUploadSize" yaml:"max-upload-size"`
	FileLocalUrl  string `mapstructure:"file-local-url" json:"fileLocalUrl" yaml:"file-local-url"`
	AllowFileType string `mapstructure:"allow-file-type" json:"allowFileType" yaml:"allow-file-type"`
}

func (f *File) AllowFileTypeList() []string {
	return strings.Split(f.AllowFileType, "|")
}
