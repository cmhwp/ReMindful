package config

//viper是啥
//viper是一个配置文件解析库，支持多种配置文件格式，如yaml、json、toml等。
import (
	"time"

	"github.com/spf13/viper"
)

// 初始化配置
type Config struct {
	Server ServerConfig `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis RedisConfig `mapstructure:"redis"`
	JWT JWTConfig `mapstructure:"jwt"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
	Mode string `mapstructure:"mode"` // 运行模式
}

type DatabaseConfig struct {
	Host string `mapstructure:"host"`
	Port int `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	Charset string `mapstructure:"charset"`
}

type RedisConfig struct {
	Host string `mapstructure:"host"`
	Port int `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB int `mapstructure:"db"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	Expiration time.Duration `mapstructure:"expiration"`
}

//加载配置
func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)//设置配置文件路径
	viper.AutomaticEnv()//自动读取环境变量

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {	
		return nil, err
	}

	return &config, nil
}
