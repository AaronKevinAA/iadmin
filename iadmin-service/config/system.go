package config
type System struct {
	Addr          int    `mapstructure:"addr" json:"addr" yaml:"addr"`                              // 端口值
	UseMultipoint bool   `mapstructure:"use-multipoint" json:"useMultipoint" yaml:"use-multipoint"` // 多点登录拦截
}
