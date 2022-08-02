package config

/**
 * @Com www.github.com/sixstaredu
 * @Author 六星教育-shineyork老师
 */

type App struct {
	Address   string `yaml:"address"`
	ServerName string `mapstructure:"server_name"`
}
