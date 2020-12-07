// Copyright (c) 2020. pkg Inc. All rights reserved.
// Author bozz
// Create Time 2020/12/7

package config

import (
	"time"
)

type Provider interface {
	Get(key string) interface{}
	GetBool(key string) bool
	GetInt(key string) int
	GetInt64(key string) int64
	GetIntSlice(key string) []int
	GetFloat64(key string) float64
	GetString(key string) string
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringSlice(key string) []string
	GetTime(key string) time.Time
	GetDuration(key string) time.Duration
	IsSet(key string) bool
	AllSettings(key string) map[string]interface{}
	Set(key string, value interface{})
}



