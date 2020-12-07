// Copyright (c) 2020. pkg Inc. All rights reserved.
// Author bozz
// Create Time 2020/12/7

package viper


import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

const (
	defaultRelativeSourcePath = "./"
	defaultEnvPrefix          = "vip"
)

var (
	errLocateConfigFile = "载入目录%s下的配置文件%s发生错误:%s"
)

// LoadConfig load config file
func loadConfig(relativeSourcePath, configFilename string) (*viper.Viper, error) {
	v := viper.New()

	v.AutomaticEnv()
	v.SetEnvPrefix(defaultEnvPrefix)

	configFilenames := strings.Split(configFilename, ",")
	for _, configFile := range configFilenames {
		v.SetConfigFile(configFile)
	}

	if relativeSourcePath == "" {
		relativeSourcePath = defaultRelativeSourcePath
	}
	v.AddConfigPath(relativeSourcePath)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil,fmt.Errorf(errLocateConfigFile, relativeSourcePath,configFilename,  err.Error())
		} else {
			return nil,err
		}
	}

	// 监听配置文件发生变化
	v.WatchConfig()
	//v.OnConfigChange(func(e fsnotify.Event) {
	//	fmt.Println("Config file changed:", e.Name)
	//})

	return v, nil
}

// 设置默认配置
func defaultConfig(cfg *viper.Viper, defaultMap map[string]interface{}) {
	for key, val := range defaultMap {
		cfg.SetDefault(key, val)
	}
}

// 选项模式
type Options struct {
	defaultCfg map[string]interface{}
	consul     []string
	etcd       []string
}

// 选项模式-最终选项值
var defaultOptions = Options{
	defaultCfg: make(map[string]interface{}),
	consul:     []string{},
	etcd:       []string{},
}

type Option func(options *Options)

func WithDefault(data map[string]interface{}) Option {
	return func(options *Options) {
		options.defaultCfg = data
	}
}

func WithConsul(endpoint, path string) Option {
	return func(options *Options) {
		options.consul = []string{"consul", endpoint, path}
	}
}

func WithEtcd(endpoint, path string) Option {
	return func(options *Options) {
		options.etcd = []string{"etcd", endpoint, path}
	}
}

func New(config string, defaultCfg map[string]interface{}) error {
	vipCfg, err := loadConfig("", config)
	if err != nil {
		return err
	}
	defaultConfig(vipCfg, defaultCfg)

	return err
}

func NewByOption(confPath string, opts ...Option) (cfg *viper.Viper, err error) {
	// 更新默认值
	for _, o := range opts {
		o(&defaultOptions)
	}

	cfg, err = loadConfig("", confPath)
	if err != nil {
		return
	}

	// 设置默认值
	if len(defaultOptions.defaultCfg) > 0 {
		defaultConfig(cfg, defaultOptions.defaultCfg)
	}

	//  设置consul
	if len(defaultOptions.consul) == 3 {
		return NewByConsul(defaultOptions.consul[1],  defaultOptions.consul[2], "toml")
		//if os.Getenv("CONSUL_HTTP_TOKEN") == "" {
		//	err = errors.New("Not found CONSUL_HTTP_TOKEN env variable")
		//	return
		//}
		//
		//cfg.AddRemoteProvider("consul", defaultOptions.consul[1], defaultOptions.consul[2])
		//cfg.SetConfigType("toml")
		//err = cfg.ReadRemoteConfig()
		//if err != nil {
		//	return
		//}
		//
		//go func() {
		//	// 使用最新配置
		//	for {
		//		time.Sleep(5 * time.Second)
		//		err := cfg.WatchRemoteConfig()
		//		if err != nil {
		//			log.Printf("unable to read remote config: %v", err)
		//			continue
		//		}
		//	}
		//}()
	}

	// 设置etcd
	if len(defaultOptions.etcd) == 3 {
		cfg.AddRemoteProvider("etcd", defaultOptions.etcd[1], defaultOptions.etcd[2])
		cfg.SetConfigType("json")
		err = cfg.ReadRemoteConfig()
		if err != nil {
			return
		}

		go func() {
			for {
				time.Sleep(5 * time.Second)
				err := cfg.WatchRemoteConfig()
				if err != nil {
					log.Printf("unable to read remote config: %v", err)
					continue
				}
			}
		}()
	}

	return
}

// 通过consul创建
func NewByConsul(endpoint, path,configType  string) (cfg *viper.Viper, err error) {
	if os.Getenv("CONSUL_HTTP_TOKEN") == "" {
		err = errors.New("Not found CONSUL_HTTP_TOKEN env variable")
		return
	}

	err = cfg.AddRemoteProvider("consul", endpoint, path)
	if err != nil {
	 return
	}

	cfg.SetConfigType(configType)
	err = cfg.ReadRemoteConfig()
	if err != nil {
		return
	}

	go func() {
		// 使用最新配置
		for {
			time.Sleep(5 * time.Second)
			err := cfg.WatchRemoteConfig()
			if err != nil {
				log.Printf("unable to read remote config: %v", err)
				continue
			}
		}
	}()

	return
}

