package cache

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/coocood/freecache"
)

// MemoryCache 是内存缓存
var MemoryCache *freecache.Cache

func init() {
	cacheSize := 100 * 1024 * 1024
	MemoryCache = freecache.NewCache(cacheSize)
	// debug.SetGCPercent(20)
}

// GetCacheData 根据项目key获取项目的详情
func GetCacheData(key string, val interface{}) (err error) {
	if len(key) <= 0 {
		log.Println("key 参数不能为空")
		return errors.New("key 参数不能为空")
	}

	v, err := MemoryCache.Get([]byte(key))
	if err != nil {
		log.Println(err)
		return err
	}

	err = json.Unmarshal(v, val)
	if err != nil {
		log.Println(err)
		return err
	}
	MemoryCache.Set([]byte(key), v, 60*60)

	return nil
}

// SetCacheData 设置内存缓存
func SetCacheData(key string, val interface{}) (err error) {
	if len(key) <= 0 {
		log.Println("key 参数不能为空")
		return errors.New("key 参数不能为空")
	}

	byts, err := json.Marshal(val)
	if err != nil {
		return err
	}

	err = MemoryCache.Set([]byte(key), byts, 60*60)
	if err != nil {
		return err
	}

	return nil
}

// DelCacheData 删除内存缓存
func DelCacheData(key string) (b bool) {
	if len(key) <= 0 {
		log.Println("key 参数不能为空")
		return false
	}

	affect := MemoryCache.Del([]byte(key))

	return affect
}
