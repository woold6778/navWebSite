package util

import (
	"time"

	"github.com/patrickmn/go-cache"
)

// 创建一个新的缓存实例，其中每项缓存默认保持5分钟，每10分钟清理一次过期项
var C = cache.New(5*time.Minute, 10*time.Minute)
