package config

import (
	"strings"
	"sync"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type LLMConfig struct {
	Provider string `mapstructure:"provider"` // 模型提供方，如 deepseek, openai
	APIKey   string `mapstructure:"api_key"`  // API 密钥
}

type Config struct {
	DoubanCookie string
	OutputDir    string
	LLM          LLMConfig `mapstructure:"llm"` // 新增
}

var (
	once   sync.Once
	config *Config
)

// Load 加载配置，优先级：环境变量 > .env 文件 > 默认值
// 环境变量支持大写形式（DEEPSEEK_key）或 BOOKNOTE_ 前缀（BOOKNOTE_DEEPSEEK_key）
func Load() *Config {
	once.Do(func() {

		// 1. 用 godotenv 加载 .env 文件（可靠解析复杂值）
		_ = godotenv.Load() // 默认当前目录 .env，也支持多路径，失败不报错

		// 2. 用 viper 读取环境变量（已包含 godotenv 注入的变量）
		v := viper.New()
		v.AutomaticEnv()
		v.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // 支持点分隔键

		// 3. 显式绑定（可选，但推荐）
		_ = v.BindEnv("douban_cookie", "DOUBAN_COOKIE")
		_ = v.BindEnv("output_dir", "OUTPUT_DIR")
		_ = v.BindEnv("llm.provider", "LLM_PROVIDER")
		_ = v.BindEnv("llm.api_key", "LLM_API_KEY", "DEEPSEEK_API_KEY")

		// 4. 默认值
		v.SetDefault("output_dir", ".")
		v.SetDefault("llm.provider", "deepseek")
		v.SetDefault("llm.api_key", "")

		config = &Config{
			DoubanCookie: v.GetString("douban_cookie"),
			OutputDir:    v.GetString("output_dir"),
			LLM: LLMConfig{
				Provider: v.GetString("llm.provider"),
				APIKey:   v.GetString("llm.api_key"),
			},
		}

	})
	return config
}

// 可选：支持从命令行标志绑定（在 main 中调用）
func BindFlags(v *viper.Viper, flags interface{}) {
	// 此函数预留，具体实现由 main 完成
}
