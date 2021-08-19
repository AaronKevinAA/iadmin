package config

import "strings"

type File struct {
	MaxUploadSize   int64  `mapstructure:"max-upload-size" json:"maxUploadSize" yaml:"max-upload-size"`
	FileLocalDir    string `mapstructure:"file-local-dir" json:"fileLocalDir" yaml:"file-local-dir"`
	AllowFileType   string `mapstructure:"allow-file-type" json:"allowFileType" yaml:"allow-file-type"`
	MaxDownloadSize int64  `mapstructure:"max-download-size" json:"maxDownloadSize" yaml:"max-download-size"`
}

func (f *File) AllowFileTypeList() []string {
	return strings.Split(f.AllowFileType, "|")
}
