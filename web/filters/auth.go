package filters

import (
	"log"
	"time"

	"github.com/caddyserver/caddy/v2/web/cache"
	"github.com/caddyserver/caddy/v2/web/model"

	"github.com/kataras/iris/v12"
)

// AuthHandle 授权中间件
func AuthHandle(ctx iris.Context) {
	// 判断用户是否登录
	var sid = ctx.GetCookie("sid")
	if len(sid) <= 0 {
		// 重定向到首页
		if ctx.IsAjax() {
			ctx.JSON(model.ResponseData{State: false, Message: "授权已过期或者无效,请登录", Data: "/", HTTPCode: 301})
			ctx.EndRequest()
			return
		}

		ctx.Redirect("/")
		ctx.EndRequest()
		return
	}

	var userSession = model.UserSession{}
	cache.GetCacheData(sid, &userSession)

	var key = "sid_" + userSession.Sid
	if key == sid && int(time.Now().Sub(userSession.CreateTime).Seconds()) < userSession.ExpireTime {
		// fmt.Println("当前登录有效")
		ctx.Values().Set("user", userSession)
		ctx.NextOrNotFound()
		return
	}

	log.Println("登录已过期")
	if ctx.IsAjax() {
		ctx.JSON(model.ResponseData{State: false, Message: "授权已过期或者无效,请登录", Data: "/", HTTPCode: 301})
		ctx.EndRequest()
		return
	}

	ctx.Redirect("/")
	ctx.EndRequest()
}
