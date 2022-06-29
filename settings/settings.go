package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"

	"github.com/spf13/viper"
)

// 全局变量 保存程序的所有配置信息
var Conf = new(Config)

type Config struct {
	*AppConfig   `mapstructure:"app"`
	*LogConfig   `mapstructure:"log"`
	*MySQLConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
}
type AppConfig struct {
	Name      string `mapstructure:"name"`
	Mode      string `mapstructure:"mode"`
	Version   string `mapstructure:"version"`
	StartTime string `mapstructure:"start_time"`
	MachineID int64  `mapstructure:"machine_id"`
	Port      int    `mapstructure:"port"`
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
	DbName       string `mapstructure:"dbname"`
	Port         int    `mapstructure:"port"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

func Init(filename string) (err error) {
	filePath := fmt.Sprintf("./conf/%s", filename)
	viper.SetConfigFile(filePath)

	//viper.SetConfigName("config") // 指定配置文件名称-不要带后缀
	//viper.SetConfigType("yaml")		// 指定配置文件类型-专门用于从远程获取配置信息时指定配置文件类型
	//viper.AddConfigPath(".") // 指定配置文件路径
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("viper.ReadInConfig failed err:%v \n", err)
		return
		//panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// 把读取到的配置信息反序列化到conf变量中
	err2 := viper.Unmarshal(Conf)
	if err2 != nil {
		fmt.Printf("viper Unmarshal failed, err:%v\n", err2)
	}

	// 监控配置文件
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Printf("配置文件：%s 修改了...\n", in.Name)
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper Unmarshal failed, err:%v\n", err2)
		}
	})
	return
}
