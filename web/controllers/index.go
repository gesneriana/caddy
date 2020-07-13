package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/caddyserver/caddy/v2/web/cache"
	"github.com/caddyserver/caddy/v2/web/filters"
	"github.com/caddyserver/caddy/v2/web/model"
	"github.com/caddyserver/caddy/v2/web/model/request"

	"net/http"

	"github.com/kataras/iris/v12"
	uuid "github.com/satori/go.uuid"
)

// RegisterIndexController 注册控制器
func RegisterIndexController(app *iris.Application) {
	// 首页, 登录页
	app.Get("/", func(ctx iris.Context) {
		// 判断用户是否登录
		var sid = ctx.GetCookie("sid")
		if len(sid) > 0 {
			var userSession = model.UserSession{}

			if userSession.Sid == sid && int(time.Now().Sub(userSession.CreateTime).Seconds()) < userSession.ExpireTime {
				fmt.Println("当前用户已登录")
				ctx.Redirect("/home/index")
				return
			}
		}

		if err := ctx.View("login.html"); err != nil {
			ctx.Application().Logger().Info(err.Error())
		}
	})

	// 网站图标
	app.Get("/favicon.ico", func(ctx iris.Context) {
		ctx.Redirect("/wwwroot/favicon.ico", 301)
	})

	// 登录api
	app.Post("/login", func(ctx iris.Context) {
		var user model.User
		var username = ctx.FormValue("username")
		var pwd = ctx.FormValue("password")

		user.UserName = username
		user.Password = pwd
		user.ID = 1

		jsondata, err := json.Marshal(user)
		if err != nil {
			log.Println(err)
		}
		if username == "admin" && pwd == "admin666" {
			// 将sid写入Redis或者数据库临时缓存
			var sid = uuid.NewV4().String()

			var userSession = model.UserSession{
				Sid:        sid,
				UserID:     user.ID,
				UserName:   user.UserName,
				CreateTime: time.Now(),
				ExpireTime: 60 * 60 * 24,
			}

			var key = "sid_" + userSession.Sid
			cache.SetCacheData(key, userSession)

			ctx.SetCookie(&http.Cookie{
				Name:     "sid",
				Value:    key,
				Path:     "/",
				HttpOnly: true,
			})
		}

		ctx.WriteString(string(jsondata))
	})

	// 需要验证cookie的授权, 分组action
	app.PartyFunc("/home", func(p iris.Party) {
		p.Use(filters.AuthHandle)

		// 首页, 显示视图静态html页面
		p.Get("/index", func(ctx iris.Context) {
			err := ctx.View("index.html")
			if err != nil {
				log.Println(err)
			}
		})
	})

	// 需要验证cookie的授权, 分组action
	app.PartyFunc("/caddy", func(p iris.Party) {
		p.Use(filters.AuthHandle)

		// 读取Caddyfile配置
		app.Get("/caddy_config", func(ctx iris.Context) {
			b, err := ioutil.ReadFile("./Caddyfile")

			if err != nil {
				fmt.Println(err.Error())
				ctx.JSON(model.ResponseData{State: false, Message: "获取Caddy配置失败", Error: err.Error(), HTTPCode: 500})
				return
			}

			ctx.JSON(model.ResponseData{State: true, Message: "获取Caddy配置成功", HTTPCode: 200, Data: string(b)})
		})

		// 保存Caddyfile配置
		app.Post("/caddy_config", func(ctx iris.Context) {
			var config = request.CaddyFileRequest{}
			var err = ctx.ReadForm(&config)
			if err != nil {
				ctx.JSON(model.ResponseData{State: false, Message: "设置Caddy配置失败", Error: err.Error(), HTTPCode: 500})
				return
			}

			f, err := os.OpenFile("./Caddyfile", os.O_RDWR|os.O_TRUNC, 0600)
			defer f.Close()
			if err != nil {
				fmt.Println(err.Error())
				ctx.JSON(model.ResponseData{State: false, Message: "设置Caddy配置失败", Error: err.Error(), HTTPCode: 500})
				return
			}
			_, err = f.Write([]byte(config.Caddy))
			if err != nil {
				fmt.Println(err.Error())
				ctx.JSON(model.ResponseData{State: false, Message: "设置Caddy配置失败", Error: err.Error(), HTTPCode: 500})
				return
			}

			ctx.JSON(model.ResponseData{State: true, Message: "设置Caddy配置成功", HTTPCode: 200, Data: config.Caddy})
		})

		// 获取Json格式的caddy配置
		app.Get("/json_config", func(ctx iris.Context) {
			resp, err := http.Get("http://127.0.0.1:2019/config")
			if err != nil {
				ctx.JSON(model.ResponseData{State: false, Message: "获取Caddy Json配置失败", Error: err.Error(), HTTPCode: 500})
				return
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)

			ctx.JSON(model.ResponseData{State: true, Message: "获取Caddy Json配置成功", HTTPCode: resp.StatusCode, Data: string(body)})
		})

		// 获取Json格式的caddy配置
		app.Post("/json_config", func(ctx iris.Context) {
			var config = model.CaddyJSONConfigModel{}
			var err = ctx.ReadJSON(&config)
			if err != nil {
				ctx.JSON(model.ResponseData{State: false, Message: "获取Caddy Json配置失败", Error: err.Error(), HTTPCode: 500})
				return
			}

			var bts, _ = json.Marshal(config)
			fmt.Println(string(bts))
			resp, err := http.Post("http://127.0.0.1:2019/load", "application/json", bytes.NewReader(bts))
			if err != nil {
				ctx.JSON(model.ResponseData{State: false, Message: "获取Caddy Json配置失败", Error: err.Error(), HTTPCode: 500})
				return
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)

			ctx.JSON(model.ResponseData{State: true, Message: "获取Caddy Json配置成功", HTTPCode: resp.StatusCode, Data: string(body)})
		})

		// 获取Json格式的caddy配置
		app.Get("/site_list", func(ctx iris.Context) {
			resp, err := http.Get("http://127.0.0.1:2019/config")
			if err != nil {
				ctx.JSON(model.ResponseData{State: false, Message: "获取Caddy Json配置失败", Error: err.Error(), HTTPCode: 500})
				return
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)

			var caddyConfig = &model.CaddyJSONConfigModel{}
			err = json.Unmarshal(body, caddyConfig)
			if err != nil {
				ctx.JSON(model.ResponseData{State: false, Message: "获取Caddy Json配置失败", Error: err.Error(), HTTPCode: 500})
				return
			}

			var routes = caddyConfig.Apps.HTTP.Servers.Srv0.Routes
			data, err := json.Marshal(routes)
			if err != nil {
				ctx.JSON(model.ResponseData{State: false, Message: "获取Caddy Json配置失败", Error: err.Error(), HTTPCode: 500})
				return
			}

			ctx.JSON(model.ResponseData{State: true, Message: "获取Caddy Json配置成功", HTTPCode: 200, Data: string(data)})
		})
	})

}
