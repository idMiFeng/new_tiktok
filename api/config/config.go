package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Conf 定义全局的变量
var Conf = new(SrvConfig)

type SrvConfig struct {
	Name          string `mapstructure:"name"`
	Mode          string `mapstructure:"mode"`
	IP            string `mapstructure:"ip"`
	Port          int    `mapstructure:"port"`
	Address       string `mapstructure:"address"`
	JwtSecret     string `mapstructure:"jwt_secret"`
	*LogConfig    `mapstructure:"log"`
	*MySQLConfig  `mapstructure:"mysql"`
	*ConsulConfig `mapstructure:"consul"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DB           string `mapstructure:"dbname"`
	Port         int    `mapstructure:"port"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}
type ConsulConfig struct {
	Addr string `mapstructure:"addr"`
}

// Init 整个服务配置文件初始化的方法
func Init(filePath string) (err error) {
	// 方式1：直接指定配置文件路径（相对路径或者绝对路径）
	// 相对路径：相对执行的可执行文件的相对路径
	// viper.SetConfigFile("./conf/config.yaml")
	viper.SetConfigFile(filePath)

	err = viper.ReadInConfig() // 读取配置信息
	if err != nil {
		// 读取配置信息失败
		fmt.Printf("viper.ReadInConfig failed, err:%v\n", err)
		return
	}
	// 如果使用的是 viper.GetXxx()方式使用配置的话，就无须下面的操作

	// 把读取到的配置信息反序列化到 Conf 变量中
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
	}

	viper.WatchConfig() // 配置文件监听
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了...")
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
		}
	})
	return
}
