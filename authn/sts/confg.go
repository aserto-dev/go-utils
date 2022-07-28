package sts

import "github.com/aserto-dev/go-utils/grpcclient"

type Config struct {
	Client grpcclient.Config `json:"client"`
	Cache  CacheOptions      `json:"cache"`
}

type CacheOptions struct {
	InvalidationTimeSeconds int `json:"invalidation_time_seconds"`
	SizeMB                  int `json:"size_mb"`
}
