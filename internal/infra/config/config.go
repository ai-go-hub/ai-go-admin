package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

// 应用配置
type App struct {
	Name string
}

// 服务配置
type Server struct {
	Port int
}

// 单个数据库连接配置
type DatabaseSource struct {
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string `mapstructure:"dbname"`
	SSLMode         string `mapstructure:"sslmode"`
	Timezone        string `mapstructure:"timezone"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

// 读库配置，额外包含 enabled 字段控制是否启用读库
type ReadConfig struct {
	Enabled        bool
	DatabaseSource `mapstructure:",squash"`
}

// 令牌配置
type TokenConfig struct {
	Driver string
}

// 图标元数据配置
type IconConfig struct {
	EnName string `mapstructure:"en_name"`
	CnName string `mapstructure:"cn_name"`
}

// 验证码开关配置，按接口控制是否启用人机验证
type CaptchaSwitches struct {
	AdminLogin bool `mapstructure:"admin_login"`
}

// 验证码配置
type CaptchaConfig struct {
	Switches        CaptchaSwitches `mapstructure:"switches"`
	Elements        []string        `mapstructure:"elements"`
	Length          int             `mapstructure:"length"`
	ConfusionLength int             `mapstructure:"confusion_length"`
	Expire          int             `mapstructure:"expire"`
	FontPath        string          `mapstructure:"font_path"`
	BackgroundDir   string          `mapstructure:"background_dir"`
	IconDir         string          `mapstructure:"icon_dir"`
	Icons           []IconConfig    `mapstructure:"icons"`
}

// 跨域配置
type CORSConfig struct {
	Origins []string
	Methods []string
	Headers []string
	MaxAge  int `mapstructure:"max_age"`
}

// 数据库配置
type DatabaseConfig struct {
	Type   string
	Prefix string
	Write  DatabaseSource
	Read   ReadConfig
}

// 聚合配置
type Config struct {
	App      App
	Server   Server
	Token    TokenConfig
	CORS     CORSConfig
	Captcha  CaptchaConfig
	Database DatabaseConfig
}

var (
	vip    *viper.Viper
	config = new(Config)
	mu     sync.RWMutex
)

// 初始化配置系统
func Init() error {
	if vip == nil {
		vip = viper.New()

		// 读取 config 目录下所有 yaml 文件
		files, err := filepath.Glob(filepath.Join("./config", "*.yaml"))
		if err != nil {
			return fmt.Errorf("get config files: %w", err)
		}

		for _, file := range files {
			vip.SetConfigFile(file)
			if err := vip.MergeInConfig(); err != nil {
				return fmt.Errorf("merge config file %s: %w", file, err)
			}
		}

		// 加载根目录 .env.yaml 文件（覆盖从 config 文件夹读取的配置）
		if _, err := os.Stat(".env.yaml"); err == nil {
			vip.SetConfigFile(".env.yaml")
			if err := vip.MergeInConfig(); err != nil {
				return fmt.Errorf("merge config file .env.yaml: %w", err)
			}
		}

		// 支持环境变量覆盖
		vip.AutomaticEnv()

		// 解析到结构体
		if err := vip.Unmarshal(config); err != nil {
			return fmt.Errorf("unmarshal config: %w", err)
		}

		// 监听配置改变
		vip.OnConfigChange(func(e fsnotify.Event) {
			mu.Lock()
			defer mu.Unlock()
			if err := vip.Unmarshal(config, func(c *mapstructure.DecoderConfig) {
				c.ZeroFields = true
			}); err != nil {
				log.Printf("config reload failed: %v", err)
			}
		})
		vip.WatchConfig()
	}
	return nil
}

// 返回全局配置副本
func Get() Config {
	mu.RLock()
	defer mu.RUnlock()
	return *config
}

// 获取 Viper 全局单例
func Viper() *viper.Viper {
	return vip
}

// 是否已经初始化
func Initialized() bool {
	return vip != nil
}
