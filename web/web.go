package web

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/caddyserver/caddy/v2/web/common"
	"github.com/caddyserver/caddy/v2/web/controllers"
	"github.com/caddyserver/caddy/v2/web/sync"

	"github.com/kataras/iris/v12"
)

func createWebAppDir() bool {
	_, err := os.Stat("./filebrowser")
	if err != nil {
		log.Printf("[Error] ./filebrowser dir not exist: %v\n", err)
		return false
	}

	_, err = os.Stat("./filebrowser/webapp")
	if err != nil {
		if os.IsNotExist(err) {
			//创建目录
			dir, _ := os.Executable()
			exPath := filepath.Dir(dir)
			if err := os.Mkdir(exPath+"/filebrowser/webapp", os.ModePerm); err != nil {
				fmt.Println(err)
				return false
			}
			return true
		}
		fmt.Println(err)
		return false
	}

	return true
}

// 初始化filebrowser模块
func initLinuxFileBrowser() {
	if createWebAppDir() == false {
		return
	}

	_, err := os.Stat("./filebrowser/filebrowser.db")
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
		shell := `cd ./filebrowser/webapp && ../filebrowser -d ../filebrowser.db`
		common.RunShellCommand(shell)
	}()

}

func initWindowsFileBrowser() {
	if createWebAppDir() == false {
		return
	}

	_, err := os.Stat("./filebrowser/filebrowser.db")
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
		shell := `cd filebrowser & cd webapp & ..\\filebrowser.exe -d ..\\filebrowser.db`
		common.RunShellCommand(shell)
	}()

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
	sync.StartFileSync()
	controllers.InitFileBrowserRoutes()

	app := iris.Default()
	// 注册视图引擎, 发布release版本Reload参数需要改为false
	app.RegisterView(iris.HTML("./wwwroot/view", ".html").Reload(true))

	// 注册视图文件目录, 不需要再为每个视图注册路由, 可以使用vue.js请求api获取数据
	app.HandleDir("/view", "./wwwroot/view")
	app.HandleDir("/wwwroot", "./wwwroot")

	controllers.RegisterIrisWebActionHandle(app)

	app.Run(iris.Addr(":2020"))
}
