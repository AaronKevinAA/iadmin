package config

type Excel struct {
	ExcelStoreDir string `mapstructure:"excel-store-dir" json:"excelStoreDir" yaml:"excel-store-dir"` // 存放excel的相对路径
}
