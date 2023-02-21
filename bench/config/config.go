package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	BaseURL           string `mapstructure:"base_url"`            // ベンチマーク先のURL
	RequestTimeout    int    `mapstructure:"request_timeout"`     // リクエストタイムアウト時間(ミリ秒)
	InitializeTimeout int    `mapstructure:"initialize_timeout"`  // 初期化タイムアウト時間(ミリ秒)
	ExitStatusOnFail  bool   `mapstructure:"exit_status_on_fail"` // ベンチマーク失敗時にexit statusを1にするかどうか
}

func init() {
	pflag.String("base_url", "http://localhost:8080", "ベンチマーク先のURL")
	pflag.Int("request_timeout", 1000, "リクエストタイムアウト時間(ミリ秒)")
	pflag.Int("initialize_timeout", 30000, "初期化タイムアウト時間(ミリ秒)")
	pflag.Bool("exit_status_on_fail", false, "ベンチマーク失敗時にexit statusを1にするかどうか")
}

func NewConfig() (c *Config, err error) {
	pflag.Parse()
	err = viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		return
	}
	viper.AutomaticEnv()
	err = viper.Unmarshal(&c)
	return
}
