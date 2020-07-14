let sitelistTemplate = Vue.extend({
    template: `<div>
    <div class='row' style='line-height:60px;'>
        <div class='col-lg-1 text-left h4' style='line-height:60px;'>网站</div>
        <div class='col-lg-9 text-left h6' style='line-height:60px;'>JSON API 接口只在服务重启之前有效,重启后读取caddyfile配置,caddyfile是永久配置,json api是临时配置</div>
        <div class='col-lg-1 text-right'><button class='btn btn-primary' @click='addSiteConfig' >添加站点</button></div>
    </div>
    <div class="row">
        <div class='col-lg-11'>
        <table class="table table-bordered table-striped" id="sitelist-table">
            <tr>
                <th style="" data-field="handle"><div class="th-inner ">处理程序(handle)<a href="https://caddyserver.com/docs/json/apps/http/servers/routes/handle/" target="_blank">帮助</a></div><div class="fht-cell"></div></th>
                <th style="" data-field="upsteams"><div class="th-inner ">后端主机(upsteams)<a href="https://caddyserver.com/docs/json/apps/http/servers/routes/handle/reverse_proxy/upstreams/" target="_blank">帮助</a></div><div class="fht-cell"></div></th>
                <th style="" data-field="match"><div class="th-inner ">匹配类型<a href="https://caddyserver.com/docs/json/apps/http/servers/routes/match/" target="_blank">帮助</a></div><div class="fht-cell"></div></th>
                <th style="" data-field="match-value"><div class="th-inner ">匹配规则</div><div class="fht-cell"></div></th>
                <th style="text-align: center; " data-field="index"><div class="th-inner ">编辑</div><div class="fht-cell"></div></th>
            </tr>
            <tr v-for="(item, index) in caddyRoutes">
                <td>
                    <span v-for="(h1, index) in item.handle" v-bind:title="h1.handler">
                        <template v-for="(h2, index) in h1.routes">
                            <template v-for="(h3, index) in h2.handle">
                                <span>{{h3.handler}}</span>
                            </template>
                        </template>
                    </span>
                </td>
                <td>
                    <template v-for="(u1, index) in item.handle">
                        <template v-for="(h2, index) in u1.routes">
                            <template v-for="(h3, index) in h2.handle">
                                <template v-for="(h4, index) in h3.upstreams">
                                    <span>{{h4.dial}}</span>
                                </template>
                            </template>
                        </template>
                    </template>
                </td>
                <td>
                    <template v-for="(m, index) in item.match">
                        <template v-for="(mtype, name) in m">
                            <span>{{name}}</span>
                        </template>
                    </template>
                </td>
                <td>
                    <template v-for="(m, index) in item.match">
                        <template v-for="(mtype, name) in m">
                            <template v-for="(mvalue, mvalue_index) in mtype">
                                <span>{{mvalue}}</span>&nbsp;
                            </template>
                        </template>
                    </template>
                </td>
                <td class="text-center">
                    <template v-for="(m, index) in item.match">
                        <template v-for="(mtype, name) in m">
                            <button v-if="name==='host'" class='btn btn-primary' v-bind:title="index" @click="editCaddySiteConfig(index)">编辑</button>
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
            this.$router.push({ path: '/newproject' })
        },
        editCaddySiteConfig: function (index) {
            // console.log(this.caddyRoutes[index]);
            // 打开编辑页面
            this.$router.push({
                name: 'editCaddySiteListConfig',
                query:
                {
                    index: index,
                    caddyConfig: JSON.stringify(this.caddyConfig)
                }
            })
        }
    },
    mounted: function () {
        var _this = this;
        $.ajax({
            type: "get",
            url: "/site_list",
            datatype: 'json',
            success: function (resp) {
                if (resp.code == 200 && resp.data != null && resp.data != "null" && resp.data.length > 0) {
                    _this.caddyConfig = JSON.parse(resp.data);
                    _this.caddyRoutes = _this.caddyConfig.apps.http.servers.srv0.routes;
                    // console.log(_this.caddyRoutes);
                }
            }
        });
    }
})