package core

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/slyrx/gin_exam_system/server/others/core/internal"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	"github.com/slyrx/gin_exam_system/server/others/global"
)

func Viper(path ...string) *viper.Viper {
	var config string

	if len(path) == 0 { // 如果没有传入路径参数
		flag.StringVar(&config, "c", "", "choose config file.")
		flag.Parse()      // 解析命令行参数
		if config == "" { // 判断命令行参数是否为空
			if configEnv := os.Getenv(internal.ConfigEnv); configEnv == "" { // 判断 internal.ConfigEnv 常量存储的环境变量是否为空
				switch gin.Mode() { // 根据 gin 的模式选择配置文件
				case gin.DebugMode:
					config = internal.ConfigDefaultFile // 使用默认配置文件
					fmt.Printf("您正在使用gin模式的%s环境名称,config的路径为%s\n", gin.Mode(), internal.ConfigDefaultFile)
				case gin.ReleaseMode:
					config = internal.ConfigReleaseFile // 使用发布配置文件
					fmt.Printf("您正在使用gin模式的%s环境名称,config的路径为%s\n", gin.Mode(), internal.ConfigReleaseFile)
				case gin.TestMode:
					config = internal.ConfigTestFile // 使用测试配置文件
					fmt.Printf("您正在使用gin模式的%s环境名称,config的路径为%s\n", gin.Mode(), internal.ConfigTestFile)
				}
			} else { // internal.ConfigEnv 常量存储的环境变量不为空，将值赋值于config
				config = configEnv
				fmt.Printf("您正在使用%s环境变量,config的路径为%s\n", internal.ConfigEnv, config)
			}
		} else { // 命令行参数不为空，将值赋值于config
			fmt.Printf("您正在使用命令行的-c参数传递的值,config的路径为%s\n", config)
		}
	} else { // 函数传递的可变参数的第一个值赋值于config
		config = path[0]
		fmt.Printf("您正在使用func Viper()传递的值,config的路径为%s\n", config)
	}

	v := viper.New()        // 创建新的 viper 实例
	v.SetConfigFile(config) // 设置配置文件路径
	v.SetConfigType("yaml") // 设置配置文件类型为 YAML
	err := v.ReadInConfig() // 读取配置文件
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err)) // 如果读取失败，抛出致命错误
	}
	v.WatchConfig() // 监控配置文件变化

	v.OnConfigChange(func(e fsnotify.Event) { // 当配置文件变化时触发
		fmt.Println("config file changed:", e.Name)
		if err = v.Unmarshal(&global.GES_CONFIG); err != nil { // 反序列化配置到全局变量
			fmt.Println(err)
		}
	})
	if err = v.Unmarshal(&global.GES_CONFIG); err != nil { // 初始反序列化配置到全局变量
		panic(err)
	}

	// root 适配性 根据root位置去找到对应迁移位置,保证root路径有效
	global.GES_CONFIG.AutoCode.Root, _ = filepath.Abs("..")
	return v
}
