package request

type Login struct {
	Phone string `json:"phone"`
	Password string `json:"password"`
	Captcha string `json:"captcha"`
	CaptchaId string `json:"captchaId"`
}

type Register struct {
	Phone string `json:"phone"`
	Password string `json:"password"`
	RealName string `json:"realName"`
	Captcha string `json:"captcha"`
	CaptchaId string `json:"captchaId"`
}

type SysUserListSearch struct {
	Pagination Pagination `gorm:"embedded"`
	Phone string		`json:"phone"`
	RealName string		`json:"realName"`
	CreatedAt []int64  `json:"createdAt"`
}