package config

import (
	"nav-web-site/util/log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// ConfigStruct 包含应用程序的所有配置
var Config ConfigStruct // 全局配置变量

// ConfigStruct 是应用程序的顶级配置结构
type ConfigStruct struct {
	Base    Base
	MySQL   MySQLConfig
	Redis   RedisConfig
	BaseUrl BaseUrlConfig `mapstructure:"base_url"`
	Tasks   []TaskConfig  `yaml:"tasks"`
}

type Base struct {
	Debug           bool
	SiteDomain      string   `mapstructure:"site_domain"`
	BackendFilePath string   `mapstructure:"backend_file_path"` //后台管理的前端程序的路径
	AllowOrigins    []string //允许的前端域名
}

type MySQLConfig struct {
	Host            string
	Port            int
	Username        string
	Password        string
	Database        string
	TablePrefix     string `mapstructure:"table_prefix"`      // 注意 mapstructure 标签// 表前缀字段
	MaxOpenConns    int    `mapstructure:"max_open_conns"`    // 最大打开连接数
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`    // 最大空闲连接数
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"` // 连接的最大生命周期（分钟）
}
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}
type BaseUrlConfig struct {
	AdPicUrl string `mapstructure:"ad_pic_url"`
}

type TaskConfig struct {
	Type     string `yaml:"type"`
	Schedule string `yaml:"schedule"`
}

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.InfoLogger.Printf("Error reading config file: %v", err)
	}

	if err := viper.Unmarshal(&Config); err != nil {
		log.InfoLogger.Printf("Unable to decode into struct: %v", err)
	}

	// 检查 TablePrefix 是否成功加载
	if Config.MySQL.TablePrefix == "" {
		log.InfoLogger.Println("Warning: TablePrefix is empty")
	} else {
		log.InfoLogger.Printf("Loaded TablePrefix: %s", Config.MySQL.TablePrefix)
	}

	//向日志输出整个配置
	log.InfoLogger.Printf("Loaded Config: %v", Config)

	// 监控配置文件变化
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.InfoLogger.Println("Config file changed:", e.Name)
		if err := viper.Unmarshal(&Config); err != nil {
			log.InfoLogger.Printf("Unable to decode into struct: %v", err)
		}
		log.InfoLogger.Printf("Reloaded Config: %v", Config)
	})
}
