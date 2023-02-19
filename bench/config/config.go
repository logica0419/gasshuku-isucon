package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	BaseURL           string `mapstructure:"base_url"`            // ベンチマーク先のURL
	RequestTimeout    int    `mapstructure:"request_timeout"`     // リクエストタイムアウト時間(ミリ秒)
	InitializeTimeout int    `mapstructure:"initialize_timeout"`  // 初期化タイムアウト時間(ミリ秒)
	ExitStatusOnFail  bool   `mapstructure:"exit_status_on_fail"` // ベンチマーク失敗時にexit statusを1にするかどうか
	Progress          bool   `mapstructure:"progress"`            // Adminロガーへの進捗の表示
}

func init() {
	viper.SetDefault("base_url", "http://localhost:8080")
	viper.SetDefault("request_timeout", 1000)
	viper.SetDefault("initialize_timeout", 30000)
	viper.SetDefault("exit_status_on_fail", false)
	viper.SetDefault("progress", true)
}

func NewConfig() (c *Config, err error) {
	viper.AutomaticEnv()
	err = viper.Unmarshal(&c)
	return
}
