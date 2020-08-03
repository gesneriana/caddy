# Caddy Web Gui

### 使用说明

文档详情请参考上游项目: https://github.com/caddyserver/caddy

预编译的二进制文件在: https://github.com/asuna6/caddy-ui

linux下的启动的命令行参数: sudo chmod 777 ./caddy && (./caddy gui &)

windows下的启动的命令行参数: ./caddy.exe gui

caddy官方主分支的功能是一个类似于NGINX的web反向代理服务器, 同时支持自动化管理https证书

此分支的主要功能是web应用的管理, web远程管理服务监听在2020端口

可以使用shell脚本在服务器上方便的部署网站, 自动同步最新发布版本的网站

并且项目集成了filebrowser模块, 可以查看网站的日志文件等

也可以直接管理网站本身的程序文件, 后台通过shell脚本来支持启动和停止功能