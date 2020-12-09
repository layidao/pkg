// Copyright (c) 2020. pkg Inc. All rights reserved.
// Author bozz
// Create Time 2020/12/7

package viper

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/layidao/pkg/config"

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

type ConfigProvider struct {
	*viper.Viper
}

// 为配置引擎拓展方法，方便符合provider接口要求
func (c *ConfigProvider) SetTest(key string, value interface{}) {
	return
}

func New(relativeSourcePath, configFilename string, defaultCfg map[string]interface{}) (cfg config.Provider, err error) {

	var (
		v *viper.Viper
	)

	v = viper.New()
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

	// 设置默认配置
	for key, val := range defaultCfg {
		v.SetDefault(key, val)
	}

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf(errLocateConfigFile, relativeSourcePath, configFilename, err.Error())
		} else {
			return nil, err
		}
	}

	// 监听配置文件发生变化
	v.WatchConfig()
	//v.OnConfigChange(func(e fsnotify.Event) {
	//	fmt.Println("Config file changed:", e.Name)
	//})

	return &ConfigProvider{v}, nil
}

// 创建使用远程的配置
func NewRemoteProvider(provider, endpoint, path, configType string, defaultCfg map[string]interface{}) (cfg config.Provider, err error) {
	var (
		v *viper.Viper
	)

	if provider == "consul" && os.Getenv("CONSUL_HTTP_TOKEN") == "" {
		err = errors.New("Not found CONSUL_HTTP_TOKEN env variable")
		return
	}

	v = viper.New()
	err = v.AddRemoteProvider(provider, endpoint, path)
	if err != nil {
		return
	}

	// 设置默认配置
	for key, val := range defaultCfg {
		v.SetDefault(key, val)
	}

	v.SetConfigType(configType)
	err = v.ReadRemoteConfig()
	if err != nil {
		// 特别注意：此处如果consulkv的值的格式不正确，也会返回not found file 错误，真坑
		return
	}

	go func() {
		// 使用最新配置
		for {
			sleep := 10 * time.Second
			if reloadDuraton := os.Getenv("VIPER_WATCH_REMOTE_DURATION"); reloadDuraton != "" {
				t, err := strconv.Atoi(reloadDuraton)
				if err == nil && t >= 10 {
					sleep = time.Duration(t) * time.Second
				}
			}

			time.Sleep(sleep)
			err := v.WatchRemoteConfig()
			if err != nil {
				log.Printf("重新从远程配置中心获取数据失败: %s", err.Error())
				continue
			}
		}
	}()

	return &ConfigProvider{v}, err
}
