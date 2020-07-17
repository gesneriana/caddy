package web

import (
	"os"

	"github.com/caddyserver/caddy/v2/web/common"
	"github.com/caddyserver/caddy/v2/web/controllers"

	"github.com/kataras/iris/v12"
)

// WebGuiStart 启动web界面
func WebGuiStart() {
	_, err := os.Stat("../../web/wwwroot")
	if err == nil {
		common.CopyDir("../../web/wwwroot", "./wwwroot")
	}

	app := iris.Default()

	// 注册视图引擎, 发布release版本Reload参数需要改为false
	app.RegisterView(iris.HTML("./wwwroot/view", ".html").Reload(true))

	// 注册视图文件目录, 不需要再为每个视图注册路由, 可以使用vue.js请求api获取数据
	app.HandleDir("/view", "./wwwroot/view")
	app.HandleDir("/wwwroot", "./wwwroot")

	controllers.RegisterIrisWebActionHandle(app)

	app.Run(iris.Addr(":2020"))
}
