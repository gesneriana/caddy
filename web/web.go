package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/caddyserver/caddy/v2/web/common"
	"github.com/caddyserver/caddy/v2/web/controllers"
	"github.com/caddyserver/caddy/v2/web/model"

	"github.com/kataras/iris/v12"
)

// 初始化filebrowser模块
func initLinuxFileBrowser() {
	_, err := os.Stat("./filebrowser")
	if err != nil {
		log.Printf("[Error] ./filebrowser dir not exist: %v\n", err)
		return
	}

	_, err = os.Stat("./filebrowser/webapp")
	if err != nil {
		if os.IsNotExist(err) {
			//创建目录
			dir, _ := os.Executable()
			exPath := filepath.Dir(dir)
			if err := os.Mkdir(exPath+"/filebrowser/webapp", os.ModePerm); err != nil {
				fmt.Println(err)
				return
			}
		}
		return
	}

	_, err = os.Stat("./filebrowser/filebrowser.db")
	if err == nil {
		log.Printf("./filebrowser/filebrowser.db already exists")
	} else {
		command := `./filebrowser.sh .`
		cmd := exec.Command("/bin/bash", "-c", command)

		output, err := cmd.Output()
		if err != nil {
			fmt.Printf("Execute Shell:%s failed with error:%s", command, err.Error())
			return
		}
		fmt.Printf("Execute Shell:%s finished with output:\n%s", command, string(output))
	}

	go func() {
		command := `cd ./filebrowser/webapp && ../filebrowser -d ../filebrowser.db`
		cmd := exec.Command("/bin/bash", "-c", command)
		output, err := cmd.Output()
		if err != nil {
			fmt.Printf("Execute Shell:%s failed with error:%s", command, err.Error())
			return
		}
		fmt.Printf("Execute Shell:%s finished with output:\n%s", command, string(output))
	}()

}

func initWindowsFileBrowser() {
	_, err := os.Stat("./filebrowser")
	if err != nil {
		log.Printf("[Error] ./filebrowser dir not exist: %v\n", err)
		return
	}

	_, err = os.Stat("./filebrowser/filebrowser.db")
	if err == nil {
		log.Printf("./filebrowser/filebrowser.db already exists")
	} else {
		cmd := exec.Command("filebrowser.bat")

		output, err := cmd.Output()
		if err != nil {
			fmt.Printf("Execute Shell: failed with error:\n%s", err.Error())
			return
		}
		fmt.Printf("Execute Shell: finished with output:\n%s", string(output))
	}

	go func() {
		command := `.\\filebrowser\\start.bat`
		cmd := exec.Command(command)
		output, err := cmd.Output()
		if err != nil {
			fmt.Printf("Execute Shell:%s failed with error:%s", command, err.Error())
			return
		}
		fmt.Printf("Execute Shell:%s finished with output:\n%s", command, string(output))
	}()

}

func initFileBrowserRoutes() {
	respJSONConfig, err := http.Get("http://127.0.0.1:2019/config")
	if err != nil {
		log.Printf("获取Caddy Json配置失败:%v\n", err)
		return
	}
	defer respJSONConfig.Body.Close()
	bodyJSONConfig, err := ioutil.ReadAll(respJSONConfig.Body)
	var config = &model.CaddyJSONConfigModel{}
	err = json.Unmarshal(bodyJSONConfig, config)
	if err != nil {
		log.Printf("解析Caddy Json配置失败:%v\n", err)
		return
	}

	var isGuiPort = false           // 是否配置了2020端口的反向代理
	var hasFileBrowserRoute = false // 2020端口是否配置了 /filebrowser
	var filebrowserIndex = 0
	var routes = config.Apps.HTTP.Servers.Srv0.Routes
	var routePath = model.RoutePath{}
	for r1Index, r1 := range routes {
		for _, r2 := range r1.Handle[0].Routes {
			var h = r2.Handle[0]
			var up = h.Upstreams[0]

			if h.Handler == "reverse_proxy" && strings.Contains(up.Dial, "2020") {
				routePathBytes, _ := json.Marshal(r2)
				json.Unmarshal(routePathBytes, &routePath)
				filebrowserIndex = r1Index
				isGuiPort = true
			}

			if filebrowserIndex == r1Index && r2.Match != nil && len(r2.Match) > 0 {
				var paths = strings.Join(r2.Match[0].Path, " ")
				if strings.Contains(paths, "/filebrowser") {
					hasFileBrowserRoute = true // 已配置path条目 /filebrowser
				}
			}
		}
	}

	if isGuiPort && hasFileBrowserRoute == false {
		routePath.Handle[0].Upstreams[0].Dial = "127.0.0.1:2021"
		if routePath.Match == nil {
			routePath.Match = make([]model.RouteMatchPath, 0)
			routePath.Match = append(routePath.Match, model.RouteMatchPath{})
		}
		routePath.Match[0].Path = []string{"/filebrowser", "/filebrowser/*"}
		routes[filebrowserIndex].Handle[0].Routes = append(routes[filebrowserIndex].Handle[0].Routes, routePath)

		var bts, _ = json.Marshal(config)
		fmt.Println(string(bts))

		respUpload, err := http.Post("http://127.0.0.1:2019/load", "application/json", bytes.NewReader(bts))
		bodyUpload, _ := ioutil.ReadAll(respUpload.Body)
		if respUpload.StatusCode != 200 || err != nil {
			log.Printf("更新Caddy Json配置失败:%v\n", err)
			return
		}
		log.Printf("更新Caddy Json配置成功:%s\n", string(bodyUpload))
	}
}

// GuiStart 启动web界面
func GuiStart() {
	_, err := os.Stat("../../web/wwwroot")
	if err == nil {
		common.CopyDir("../../web/wwwroot", "./wwwroot")
	}

	if runtime.GOOS == "windows" {
		initWindowsFileBrowser()
	} else if runtime.GOOS == "linux" {
		initLinuxFileBrowser()
	}
	initFileBrowserRoutes()

	app := iris.Default()
	// 注册视图引擎, 发布release版本Reload参数需要改为false
	app.RegisterView(iris.HTML("./wwwroot/view", ".html").Reload(true))

	// 注册视图文件目录, 不需要再为每个视图注册路由, 可以使用vue.js请求api获取数据
	app.HandleDir("/view", "./wwwroot/view")
	app.HandleDir("/wwwroot", "./wwwroot")

	controllers.RegisterIrisWebActionHandle(app)

	app.Run(iris.Addr(":2020"))
}
