package controllers

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/caddyserver/caddy/v2/web/cache"
	"github.com/caddyserver/caddy/v2/web/common"
	"github.com/caddyserver/caddy/v2/web/controllers/actions"
	"github.com/caddyserver/caddy/v2/web/filters"
	"github.com/caddyserver/caddy/v2/web/model"
	"github.com/caddyserver/caddy/v2/web/model/request"

	"net/http"

	"github.com/kataras/iris/v12"
	uuid "github.com/satori/go.uuid"
)

// InitFileBrowserRoutes 初始化 filebrowser 模块的路由path
func InitFileBrowserRoutes() {
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
		// 解决路由优先匹配 / 所以必须将 /filebrowser放在json数组的前面, JSON API必须注意顺序, caddyfile不受影响
		var pathRoutes = routes[filebrowserIndex].Handle[0].Routes
		routes[filebrowserIndex].Handle[0].Routes = pathRoutes[1:]
		routes[filebrowserIndex].Handle[0].Routes = append(routes[filebrowserIndex].Handle[0].Routes, pathRoutes[0])

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

// RegisterIrisWebActionHandle 注册iris的路由web handle
func RegisterIrisWebActionHandle(app *iris.Application) {
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
		pwdHash := sha256.Sum256([]byte(pwd))
		userHash := sha256.Sum256([]byte(username + pwd))
		pwdHashString := base64.StdEncoding.EncodeToString(pwdHash[:])
		userHashString := base64.StdEncoding.EncodeToString(userHash[:])

		user.UserName = username
		user.Password = pwd
		user.ID = 1

		userJSONBts, err := ioutil.ReadFile("./pwd.json")
		// 成功读取配置文件才会校验密码, 否则使用默认的用户名和密码
		if err == nil {
			var userConfig = &model.CaddyUser{}
			err = json.Unmarshal(userJSONBts, userConfig)
			if err != nil {
				ctx.JSON(model.ResponseData{State: false, Message: "用户登录失败", Error: err.Error(), HTTPCode: 500})
				return
			}

			if userConfig.UserName == username && userConfig.PasswordHash == pwdHashString {
				// 将sid写入Redis或者数据库临时缓存
				var sid = uuid.NewV4().String()

				var userSession = model.UserSession{
					Sid:        sid,
					UserID:     user.ID,
					UserName:   user.UserName,
					UserHash:   userHashString,
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
				ctx.JSON(model.ResponseData{State: true, Message: "用户登录成功", HTTPCode: 200})
				return
			}

			ctx.JSON(model.ResponseData{State: false, Message: "用户登录失败", HTTPCode: 403})
			return
		}

		if username == "admin" && pwd == "admin" {
			// 将sid写入Redis或者数据库临时缓存
			var sid = uuid.NewV4().String()

			var userSession = model.UserSession{
				Sid:        sid,
				UserID:     user.ID,
				UserName:   user.UserName,
				UserHash:   userHashString,
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
			ctx.JSON(model.ResponseData{State: true, Message: "用户登录成功", HTTPCode: 200})
			return
		}

		ctx.JSON(model.ResponseData{State: false, Message: "用户登录失败", HTTPCode: 403})
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

		p.Post("/changepwd", func(ctx iris.Context) {
			username := ctx.FormValue("username")
			oldpwd := ctx.FormValue("oldpassword")
			newpwd := ctx.FormValue("newpassword")
			var sid = ctx.GetCookie("sid")

			oldpwdHash := sha256.Sum256([]byte(oldpwd))
			newpwdHash := sha256.Sum256([]byte(newpwd))
			oldpwdHashString := base64.StdEncoding.EncodeToString(oldpwdHash[:])
			newpwdHashString := base64.StdEncoding.EncodeToString(newpwdHash[:])

			var user = &model.CaddyUser{}

			_, err := os.Stat("./pwd.json")
			if err != nil {
				if !os.IsNotExist(err) {
					ctx.JSON(model.ResponseData{State: false, Message: "更新用户设置失败", Error: err.Error(), HTTPCode: 500})
					return
				}

				f, err := os.Create("./pwd.json")
				defer f.Close()
				if err != nil {
					ctx.JSON(model.ResponseData{State: false, Message: "更新用户设置失败", Error: err.Error(), HTTPCode: 500})
					return
				}

				// pwd.json不存在, 不做校验, 直接写入配置文件
				user.UserName = username
				user.PasswordHash = newpwdHashString
				bts, _ := json.Marshal(user)
				f.Write(bts)

				cache.DelCacheData(sid)
				ctx.JSON(model.ResponseData{State: true, Message: "更新用户设置成功", HTTPCode: 200, Data: string(bts)})
				return
			}

			btsPwdFile, err := ioutil.ReadFile("./pwd.json")
			err = json.Unmarshal(btsPwdFile, user)
			if err != nil {
				ctx.JSON(model.ResponseData{State: false, Message: "更新用户设置失败", Error: err.Error(), HTTPCode: 500})
				return
			}

			// 只校验密码, 不校验用户名, 因为当前用户已经是登录状态了
			if user.PasswordHash != oldpwdHashString {
				ctx.JSON(model.ResponseData{State: false, Message: "更新用户设置失败", Data: "旧密码不正确", HTTPCode: 400})
				return
			}

			user.UserName = username
			user.PasswordHash = newpwdHashString
			btsUser, _ := json.Marshal(user)
			err = ioutil.WriteFile("./pwd.json", btsUser, os.ModePerm)
			if err != nil {
				ctx.JSON(model.ResponseData{State: false, Message: "更新用户设置失败", Error: err.Error(), HTTPCode: 500})
				return
			}

			cache.DelCacheData(sid)
			ctx.JSON(model.ResponseData{State: true, Message: "更新用户设置成功", HTTPCode: 200, Data: string(btsUser)})
			return
		})

		p.Get("/token", func(ctx iris.Context) {
			u := ctx.Values().Get("user")
			if user, ok := u.(model.UserSession); ok {
				ctx.JSON(model.ResponseData{State: true, Message: "获取token成功", Data: user.UserHash, HTTPCode: 200})
				return
			}

			ctx.JSON(model.ResponseData{State: false, Message: "获取token失败", HTTPCode: 500})
		})

		p.Get("/downloadCert", func(ctx iris.Context) {
			crtPath := ctx.URLParam("crt")
			certname := ctx.URLParam("certname")
			if strings.Contains(crtPath, "/.local/share/caddy/certificates") == false {
				return
			}

			if len(certname) == 0 {
				certname = "cert"
			}

			log.Println("crtPath:" + crtPath)
			err := ctx.SendFile(crtPath, certname+".crt")
			if err != nil {
				ctx.WriteString(err.Error())
			}
		})

		p.Get("/downloadCertKey", func(ctx iris.Context) {
			keyPath := ctx.URLParam("key")
			certname := ctx.URLParam("certname")

			if strings.Contains(keyPath, "/.local/share/caddy/certificates") == false {
				return
			}

			if len(certname) == 0 {
				certname = "cert"
			}

			log.Println("keyPath:" + keyPath)
			err := ctx.SendFile(keyPath, certname+".key")
			if err != nil {
				ctx.WriteString(err.Error())
			}
		})

		p.Get("/filebrowsertoken", func(ctx iris.Context) {
			var key = "jwt"
			jwtBytes, err := cache.MemoryCache.Get([]byte(key))
			if err == nil && len(jwtBytes) > 0 {
				ctx.JSON(model.ResponseData{State: true, Message: "filebrowser服务登录成功", Data: string(jwtBytes), HTTPCode: 200})
				return
			}

			var user = &request.FileBrowserLogin{}
			user.UserName = "admin"
			user.Password = "filebrowseradmin"
			user.Recaptcha = ""

			var bts, _ = json.Marshal(user)
			fmt.Println(string(bts))
			resp, err := http.Post("http://127.0.0.1:2021/filebrowser/api/login", "application/json", bytes.NewReader(bts))
			body, _ := ioutil.ReadAll(resp.Body)
			if resp.StatusCode != 200 || err != nil {
				ctx.JSON(model.ResponseData{State: false, Message: "设置Caddy配置失败", Error: err.Error(), Data: body, HTTPCode: 500})
				return
			}

			cache.MemoryCache.Set([]byte(key), body, 60*60)
			ctx.JSON(model.ResponseData{State: true, Message: "filebrowser服务登录成功", Data: body, HTTPCode: 200})
			return
		})

		p.Post("/filebrowserpath", func(ctx iris.Context) {
			var domain = ctx.FormValue("domain")
			_, err := os.Stat("./filebrowser/webapp/" + domain)
			if err != nil {
				if os.IsNotExist(err) {
					//创建目录
					dir, _ := os.Executable()
					exPath := filepath.Dir(dir)
					if err := os.Mkdir(exPath+"/filebrowser/webapp/"+domain, os.ModePerm); err != nil {
						ctx.JSON(model.ResponseData{State: false, Message: "创建web app目录失败", Error: err.Error(), HTTPCode: 500})
						return
					}

					ctx.JSON(model.ResponseData{State: true, Message: "创建web app目录成功", Data: "/filebrowser/files/" + domain, HTTPCode: 200})
					return
				}

				ctx.JSON(model.ResponseData{State: false, Message: "创建web app目录失败", Error: err.Error(), HTTPCode: 500})
				return
			}

			ctx.JSON(model.ResponseData{State: true, Message: "web app目录已存在", Data: "/filebrowser/files/" + domain, HTTPCode: 200})
			return
		})

		p.Get("/GitSyncConfig", func(ctx iris.Context) {
			var shellConfigMap map[string]model.ShellConfig
			shellConfigMap = make(map[string]model.ShellConfig)
			err := cache.GetCacheData("SyncShellConfigList", &shellConfigMap)
			if err != nil {
				ctx.JSON(model.ResponseData{State: false, Message: "从缓存中读取配置失败", Error: err.Error(), HTTPCode: 500})
				return
			}
			bts, _ := json.Marshal(shellConfigMap)
			ctx.JSON(model.ResponseData{State: true, Message: "从缓存中读取配置成功", Data: string(bts), HTTPCode: 200})
		})

		p.Post("/GitSyncConfig", func(ctx iris.Context) {
			var shellConfig = &model.ShellConfig{}
			err := ctx.ReadForm(shellConfig)
			if err != nil || len(shellConfig.Domain) == 0 {
				ctx.JSON(model.ResponseData{State: false, Message: "读取表单内容失败", Error: err.Error(), HTTPCode: 500})
				return
			}

			if shellConfig.Interval < 10 || shellConfig.Interval > 600 {
				shellConfig.Interval = 60
			}

			_, err = os.Stat("./filebrowser/webapp/" + shellConfig.Domain)
			if err != nil {
				if os.IsNotExist(err) {
					//创建目录
					dir, _ := os.Executable()
					exPath := filepath.Dir(dir)
					if err := os.Mkdir(exPath+"/filebrowser/webapp/"+shellConfig.Domain, os.ModePerm); err != nil {
						fmt.Println(err)
						ctx.JSON(model.ResponseData{State: false, Message: "写入配置失败", Error: err.Error(), HTTPCode: 500})
						return
					}
				}
				ctx.JSON(model.ResponseData{State: false, Message: "写入配置失败", Error: err.Error(), HTTPCode: 500})
				return
			}

			// 非必填, 可以手动上传可执行程序到webapp下的 domain 目录
			if len(shellConfig.InitShell) > 0 {
				var shell = ""
				if runtime.GOOS == "windows" {
					shell = fmt.Sprintf("cd ./filebrowser/webapp/%s & ", shellConfig.Domain)
				} else if runtime.GOOS == "linux" {
					shell = fmt.Sprintf("cd ./filebrowser/webapp/%s && ", shellConfig.Domain)
				}
				result, err := common.RunShellCommand(shell + shellConfig.InitShell)
				if err != nil {
					shellConfig.IsInit = false
					fmt.Println("初始化shell执行失败:", result, err)
				} else {
					shellConfig.IsInit = true
				}
			}

			bts, _ := json.Marshal(shellConfig)
			err = ioutil.WriteFile("./filebrowser/webapp/"+shellConfig.Domain+"/shell_config.json", bts, os.ModePerm)
			if err != nil {
				ctx.JSON(model.ResponseData{State: false, Message: "写入配置失败", Error: err.Error(), HTTPCode: 500})
				return
			}

			var shellConfigMap map[string]model.ShellConfig
			shellConfigMap = make(map[string]model.ShellConfig)
			err = cache.GetCacheData("SyncShellConfigList", &shellConfigMap)
			if err == nil {
				shellConfigMap[shellConfig.Domain] = *shellConfig
				cache.SetCacheData("SyncShellConfigList", shellConfigMap)
			}

			ctx.JSON(model.ResponseData{State: true, Message: "写入配置成功", Data: string(bts), HTTPCode: 200})
		})

		p.Post("/ExecShell", actions.ExecShell)
	})

	// 需要验证cookie的授权, 分组action
	app.PartyFunc("/caddy", func(p iris.Party) {
		p.Use(filters.AuthHandle)

		// 读取Caddyfile配置
		p.Get("/caddy_config", func(ctx iris.Context) {
			b, err := ioutil.ReadFile("./Caddyfile")

			if err != nil {
				fmt.Println(err.Error())
				ctx.JSON(model.ResponseData{State: false, Message: "获取Caddy配置失败", Error: err.Error(), HTTPCode: 500})
				return
			}

			ctx.JSON(model.ResponseData{State: true, Message: "获取Caddy配置成功", HTTPCode: 200, Data: string(b)})
		})

		// 保存Caddyfile配置
		p.Post("/caddy_config", func(ctx iris.Context) {
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

			// 更新caddyfile后, gowatch不是立即生效, 每秒检查一次文件更改
			go func() {
				time.Sleep(time.Second * 3)
				InitFileBrowserRoutes()
			}()

			ctx.JSON(model.ResponseData{State: true, Message: "设置Caddy配置成功", HTTPCode: 200, Data: config.Caddy})
		})

		// 获取Json格式的caddy配置
		p.Get("/json_config", func(ctx iris.Context) {
			resp, err := http.Get("http://127.0.0.1:2019/config")
			if err != nil {
				ctx.JSON(model.ResponseData{State: false, Message: "获取Caddy Json配置失败", Error: err.Error(), HTTPCode: 500})
				return
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)

			ctx.JSON(model.ResponseData{State: true, Message: "获取Caddy Json配置成功", HTTPCode: resp.StatusCode, Data: string(body)})
		})

		// 通过Caddy JSON API管理caddy web服务器, 只在运行期间生效, 重启后失效, 暂不考虑做JSON API的持久化
		p.Post("/json_config", func(ctx iris.Context) {
			var config = model.CaddyJSONConfigModel{}
			var err = ctx.ReadJSON(&config)
			if err != nil {
				ctx.JSON(model.ResponseData{State: false, Message: "设置Caddy Json配置失败", Error: err.Error(), HTTPCode: 500})
				return
			}

			var bts, _ = json.Marshal(config)
			fmt.Println(string(bts))

			/*
				caddyConfig, err := caddyfile.FromJSON(bts)

				if err != nil {
					ctx.JSON(model.ResponseData{State: false, Message: "Json转换Caddy配置失败", Error: err.Error(), HTTPCode: 500})
					return
				}
			*/

			// 本次修改有效, 下次重启caddy后json配置会丢失, 再次读取caddy file
			resp, err := http.Post("http://127.0.0.1:2019/load", "application/json", bytes.NewReader(bts))
			body, _ := ioutil.ReadAll(resp.Body)
			if resp.StatusCode != 200 || err != nil {
				ctx.JSON(model.ResponseData{State: false, Message: "设置Caddy配置失败", Error: err.Error(), Data: body, HTTPCode: 500})
				return
			}

			// Caddyfile配置文件和json格式配置文件之间并不兼容, 暂不实现转换接口
			/*
				f, err := os.OpenFile("./Caddyfile", os.O_RDWR|os.O_TRUNC, 0600)
				defer f.Close()
				if err != nil {
					ctx.JSON(model.ResponseData{State: false, Message: "设置Caddy配置失败", Error: err.Error(), HTTPCode: 500})
					return
				}

				_, err = f.Write(caddyConfig)
				if err != nil {
					ctx.JSON(model.ResponseData{State: false, Message: "设置Caddy配置失败", Error: err.Error(), HTTPCode: 500})
					return
				}
			*/
			ctx.JSON(model.ResponseData{State: true, Message: "设置Caddy Json配置成功", HTTPCode: resp.StatusCode, Data: string(body)})
		})

		// 获取Json格式的caddy配置
		p.Get("/site_list", func(ctx iris.Context) {
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

			// var routes = caddyConfig.Apps.HTTP.Servers.Srv0.Routes
			data, err := json.Marshal(caddyConfig)
			if err != nil {
				ctx.JSON(model.ResponseData{State: false, Message: "获取Caddy Json配置失败", Error: err.Error(), HTTPCode: 500})
				return
			}

			ctx.JSON(model.ResponseData{State: true, Message: "获取Caddy Json配置成功", HTTPCode: 200, Data: string(data)})
		})

		p.Get("/certlist", func(ctx iris.Context) {
			homeDir, err := common.Home()
			if err != nil {
				ctx.JSON(model.ResponseData{State: false, Message: "获取Caddy 证书列表失败", Error: err.Error(), HTTPCode: 500})
				return
			} else if len(homeDir) == 0 {
				ctx.JSON(model.ResponseData{State: false, Message: "获取Caddy 证书列表失败", Data: "homeDir长度为0", HTTPCode: 500})
				return
			}

			var certData = model.CertData{CertList: make([]model.CertModel, 0)}
			certRootPath := homeDir + "/.local/share/caddy/certificates"
			acmeDirs, _ := ioutil.ReadDir(certRootPath)
			for _, d1 := range acmeDirs {
				domainDirs, _ := ioutil.ReadDir(certRootPath + "/" + d1.Name())
				for _, domainDir := range domainDirs {
					certDir := certRootPath + "/" + d1.Name() + "/" + domainDir.Name()
					var certModel = model.CertModel{}
					certModel.CertDir = certDir
					certModel.Domain = domainDir.Name()
					certModel.LastModifiedTime = domainDir.ModTime()
					certData.CertList = append(certData.CertList, certModel)
				}
			}

			certDataBytes, err := json.Marshal(certData)
			if err != nil {
				ctx.JSON(model.ResponseData{State: false, Message: "获取Caddy 证书列表失败", Error: err.Error(), HTTPCode: 500})
				return
			}

			ctx.JSON(model.ResponseData{State: true, Message: "获取Caddy 证书列表成功", HTTPCode: 200, Data: string(certDataBytes)})
		})
	})

}
