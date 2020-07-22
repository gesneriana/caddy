let webapplistTemplate = Vue.extend({
    template: `<div>
    <div class='row' style='line-height:60px;'>
        <div class='col-lg-2 text-left h4' style='line-height:60px;'>Web应用</div>
        <div class='col-lg-8 text-left h6' style='line-height:60px;'></div>
        <div class='col-lg-1 text-right'><button class='btn btn-primary' @click='addSiteConfig' >添加站点</button></div>
    </div>
    <div class="row">
        <div class='col-lg-11'>
        <table class="table table-bordered table-striped" id="sitelist-table">
            <tr>
                <th style="" data-field="upsteams"><div class="th-inner ">后端主机(upsteams)<a href="https://caddyserver.com/docs/json/apps/http/servers/routes/handle/reverse_proxy/upstreams/" target="_blank">帮助</a></div><div class="fht-cell"></div></th>
                <th style="" data-field="match"><div class="th-inner ">匹配类型<a href="https://caddyserver.com/docs/json/apps/http/servers/routes/match/" target="_blank">帮助</a></div><div class="fht-cell"></div></th>
                <th style="" data-field="match-value"><div class="th-inner ">匹配域名</div><div class="fht-cell"></div></th>
                <th style="text-align: center; " data-field="index"><div class="th-inner ">编辑</div><div class="fht-cell"></div></th>
            </tr>
            <tr v-for="(item, index) in caddyRoutes">
                <td>
                    <template v-for="u1 in item.handle">
                        <template v-for="h2 in u1.routes">
                            <div class="row">
                                <div class="col-lg-6">
                                    <template v-for="h3 in h2.handle">
                                        <template v-for="h4 in h3.upstreams">
                                            <span>{{h4.dial}}</span>
                                        </template>
                                    </template>
                                </div>
                                <div class="col-lg-6">
                                    <template v-for="h3 in h2.match">
                                        <template v-for="h4 in h3.path">
                                            <span>{{h4}}&nbsp;</span>
                                        </template>
                                    </template>
                                </div>
                            </div>
                        </template>
                    </template>
                </td>
                <td>
                    <template v-for="m1 in item.match">
                        <template v-for="(mtype, name) in m1">
                            <span>{{name}}</span>
                        </template>
                    </template>
                </td>
                <td>
                    <template v-for="m2 in item.match">
                        <template v-for="(mtype, name) in m2">
                            <template v-for="mvalue in mtype">
                                <span>{{mvalue}}</span>&nbsp;
                            </template>
                        </template>
                    </template>
                </td>
                <td class="text-center">
                    <template v-for="m3 in item.match">
                        <template v-for="(mtype, name) in m3">
                            <button v-if="name==='host'" class='btn btn-primary' @click="gitPull(index)">远程同步</button>
                            <button v-if="name==='host'" class='btn btn-primary' @click="uploadWebApp(index)">文件管理</button>
                            <button v-if="name==='host'" class='btn btn-primary' @click="editStartShellScripts(index)">启动脚本</button>
                            <button v-if="name==='host'" class='btn btn-primary' >停止脚本</button>
                        </template>
                    </template>
                </td>
            </tr>
        </table>
        </div>
    </div>
    </div>`,
    data: function () {
        return {
            caddyConfig: {},
            caddyRoutes: []
        }
    },
    methods: {
        addSiteConfig: function () {
            this.$router.push({
                path: '/addsite',
                query: {
                    caddyConfig: JSON.stringify(this.caddyConfig)
                }
            })
        },
        uploadWebApp: function (index) {
            var _this = this;
            var route = _this.caddyRoutes[index];
            var domain = route.match[0].host[0];
            // 请求后端创建 www.example.com 目录, 然后打开新窗口显示此文件夹
            $.ajax({
                type: "post",
                url: "/home/filebrowserpath",
                data: "domain=" + domain,
                datatype: 'json',
                success: function (resp) {
                    if (resp.code == 200 && resp.state == true) {
                        window.open(resp.data, "_blank");
                        // console.log(_this.caddyRoutes);
                    }
                }
            });
        },
        editStartShellScripts: function (index) {
            // console.log(this.caddyRoutes[index]);
            this.caddyRoutes.splice(index, 1);
            var _this = this;
            // 还需要优化filebrowser模块的caddy路由自动化配置, 将端口为2020的域名添加 path路由匹配filebrowser模块
            console.log(this.caddyConfig);
        },
    },
    mounted: function () {
        var _this = this;
        $.ajax({
            type: "get",
            url: "/caddy/site_list",
            datatype: 'json',
            success: function (resp) {
                if (resp.code == 200 && resp.data != null && resp.data != "null" && resp.data.length > 0) {
                    _this.caddyConfig = JSON.parse(resp.data);
                    _this.caddyRoutes = _this.caddyConfig.apps.http.servers.srv0.routes;
                    // console.log(_this.caddyRoutes);
                }
            }
        });

        $.ajax({
            type: "get",
            url: "/home/filebrowsertoken",
            datatype: 'json',
            success: function (resp) {
                if (resp.code == 200 && resp.data != null && resp.data != "null" && resp.data.length > 0) {
                    // console.log(_this.caddyRoutes);
                    localStorage.setItem("jwt", resp.data);
                }
            }
        });
    }
})