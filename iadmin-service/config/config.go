package config

type Server struct {
	JWT     JWT     `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Zap     Zap     `mapstructure:"zap" json:"zap" yaml:"zap"`
	Redis   Redis   `mapstructure:"redis" json:"redis" yaml:"redis"`
	Captcha Captcha `mapstructure:"captcha" json:"captcha" yaml:"captcha"`
	// gorm
	Mysql  Mysql  `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	System System `mapstructure:"system" json:"system" yaml:"system"`
	Casbin Casbin `mapstructure:"casbin" json:"casbin" yaml:"casbin"`
	File   File   `mapstructure:"file" json:"file" yaml:"file"`
	Excel  Excel  `mapstructure:"excel" json:"excel" yaml:"excel"`
}
