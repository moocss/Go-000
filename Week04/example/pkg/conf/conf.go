package conf

import (
	"fmt"
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/pkg/errors"
)

type Conf struct {
	viper *viper.Viper
}

func NewConfig(path string) *Conf {
	v := viper.New()
	v.SetConfigFile(path)

	// Config's format: "json" | "toml" | "yaml" | "yml"
	v.SetConfigType("yaml")

	// viper 解析配置文件
	if err := v.ReadInConfig(); err != nil {
		log.Printf("Using config file: %s [%s]\n", v.ConfigFileUsed(), err)
		return nil
	}

	// 读取匹配的环境变量
	v.AutomaticEnv()

	return &Conf{v}
}

// WatchConfig 监控配置文件变化并热加载程序
func (c *Conf) WatchConfig() {
	c.viper.WatchConfig()
	c.viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Config file changed: %s \n", e.Name)
	})
}

// GetFloat64 获取浮点数配置
func (c *Conf) GetFloat64(key string) float64 {
	return c.viper.GetFloat64(key)
}

// Get 获取字符串配置
func (c *Conf) Get(key string) string {
	return c.viper.GetString(key)
}

// GetInt 获取整数配置
func (c *Conf) GetInt(key string) int {
	return c.viper.GetInt(key)
}

// GetInt32 获取 int32 配置
func (c *Conf) GetInt32(key string) int32 {
	return c.viper.GetInt32(key)
}

// GetInt64 获取 int64 配置
func (c *Conf) GetInt64(key string) int64 {
	return c.viper.GetInt64(key)
}

// GetDuration 获取时间配置
func (c *Conf) GetDuration(key string) time.Duration {
	return c.viper.GetDuration(key)
}

// GetTime 查询时间配置
// 默认时间格式为 "2006-01-02 15:04:05"，conf.GetTime("FOO_BEGIN")
// 如果需要指定时间格式，则可以多传一个参数，conf.GetString("FOO_BEGIN", "2006")
//
// 配置不存在或时间格式错误返回**空时间对象**
// 使用本地时区
func (c *Conf) GetTime(key string, args ...string) time.Time {
	fmt := "2006-01-02 15:04:05"
	if len(args) == 1 {
		fmt = args[0]
	}

	t, _ := time.ParseInLocation(fmt, c.viper.GetString(key), time.Local)
	return t
}

// GetBool 获取配置布尔配置
func (c *Conf) GetBool(key string) bool {
	return c.viper.GetBool(key)
}

// Set 设置配置，仅用于测试
func (c *Conf) Set(key string, value string) {
	c.viper.Set(key, value)
}

// Sub 返回新的Viper实例，代表该实例的子节点。
func (c *Conf) Sub(key string) (*viper.Viper, error) {
	if app := c.viper.Sub(key); app != nil {
		return app, nil
	}
	return nil, errors.New(fmt.Sprintf("No found `%s` in the configuration", key))
}
