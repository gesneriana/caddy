package sync

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"runtime"
	"time"

	"github.com/caddyserver/caddy/v2/web/cache"
	"github.com/caddyserver/caddy/v2/web/common"
	"github.com/caddyserver/caddy/v2/web/model"
)

// StartFileSync 开启一个goroutine同步文件
func StartFileSync() {
	// 读取配置文件的同步脚本
	rd, err := ioutil.ReadDir("./filebrowser/webapp")
	if err != nil {
		fmt.Println(err)
		return

	}

	var configMap map[string]model.ShellConfig
	configMap = make(map[string]model.ShellConfig)
	for _, fi := range rd {
		if fi.IsDir() {
			bytes, err := ioutil.ReadFile("./filebrowser/webapp/" + fi.Name() + "/shell_config.json")
			if err != nil {
				fmt.Println(err)
				continue
			}

			var config = &model.ShellConfig{}
			err = json.Unmarshal(bytes, config)
			if err != nil {
				fmt.Println(err)
				continue
			}

			configMap[fi.Name()] = *config
		}
	}
	cache.SetCacheData("SyncShellConfigList", configMap)

	go func() {
		for {
			time.Sleep(time.Second * 10)
			var shellConfigMap map[string]model.ShellConfig
			shellConfigMap = make(map[string]model.ShellConfig)
			err := cache.GetCacheData("SyncShellConfigList", &shellConfigMap)
			if err != nil {
				return
			}

			for k, v := range shellConfigMap {
				if len(v.SyncShell) > 0 && int(time.Now().Sub(v.LastSyncTime).Seconds()) > v.Interval {
					var shell = ""
					if runtime.GOOS == "windows" {
						shell = fmt.Sprintf("cd ./filebrowser/webapp/%s & ", k)
					} else if runtime.GOOS == "linux" {
						shell = fmt.Sprintf("cd ./filebrowser/webapp/%s && ", k)
					}

					common.RunShellCommand(shell + v.SyncShell)
					v.LastSyncTime = time.Now()
					shellConfigMap[k] = v
				}
			}

			cache.SetCacheData("SyncShellConfigList", shellConfigMap)
		}
	}()
}
