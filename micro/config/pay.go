package config

/**
 * @Com www.github.com/sixstaredu
 * @Author 六星教育-shineyork老师
 */

type Pay struct {
	AppId 		 string `mapstructure:"app_id"`
	AliPublicKey string `mapstructure:"ali_public_key"`
	PrivateKey 	 string `mapstructure:"private_key"`
	NotifyURL	 string `mapstructure:"notify_url"`
}