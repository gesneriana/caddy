let editCaddySiteListConfigTemplate = Vue.extend({
    template: `<div>
    <div class='row' style='line-height:60px;'>
        <div class='col-lg-2 text-left h4' style='line-height:60px;'>编辑项目</div>
        <div class='col-lg-10 text-right' style='padding-right:2em;'></div>
    </div>
    <div>
        <form id="caddyform" action="/home/editProject" method="POST" class="form-horizontal" role="form">
            <div class="form-group">
                <label for="dial" class="control-lable col-sm-2">后端主机</label>
                <div class="col-sm-10">
                    <input id="dial" name="dial" class="form-control" type="text" v-model="dial" >
                </div>
            </div>
            <div class="form-group">
                <label class="control-lable col-lg-2">域名</label>
                <div class="col-lg-10">
                    <input name="ProjectPlatform" class="form-control" type="text" v-model="host_list" >
                </div>
            </div>

            <div class="form-group">
                <div class="col-sm-10 col-sm-offset-2">
                    <button type="button" class="btn btn-primary" @click='submit($event)'>保存</button>
                </div>
            </div>
        </form>
    </div>
        </div>
    `,
    data: function () {
        return {
            caddyConfig: {},
            caddyRoutes: [],
            index: 0,
            editRoute: {},
            dial: "",
            host_list: ""
        }
    },
    methods: {
        submit: function (event) {
            // 组装请求参数
            if (this.dial.trim().length > 0 && this.host_list.trim().length > 0) {
                var hosts = this.host_list.split(",");
                this.editRoute.match[0].host = hosts;

                this.editRoute.handle[0].routes[0].handle[0].upstreams[0].dial = this.dial;
                this.caddyRoutes[this.index] = this.editRoute;
                this.caddyConfig.apps.http.servers.srv0.routes = this.caddyRoutes;

                var _this = this;
                $.ajax({
                    type: "post",
                    url: "/json_config",
                    contentType: 'application/json',  //指定格式为json格式
                    datatype: 'json',
                    data: JSON.stringify(_this.caddyConfig),
                    success: function (resp) {
                        console.log(resp);
                        if (resp.code == 200 && resp.data != null && resp.data != "null") {
                            console.log(resp.data);
                            _this.$router.go(-1);
                        }
                    }
                });

                return;
            }

            console.warn("参数错误,后端主机和域名必须填写有效内容.");
        }
    },
    mounted: function () {
        var _this = this;
        this.index = this.$route.query.index;
        this.caddyConfig = JSON.parse(this.$route.query.caddyConfig);
        this.caddyRoutes = this.caddyConfig.apps.http.servers.srv0.routes;
        this.editRoute = this.caddyRoutes[this.index];
        this.dial = this.editRoute.handle[0].routes[0].handle[0].upstreams[0].dial;
        for (var i = 0; i < this.editRoute.match[0].host.length; i++) {
            this.host_list += this.editRoute.match[0].host[i] + ","
        }
        this.host_list = this.host_list.substring(0, this.host_list.length - 1);
        // console.log(this.$route.query);
    }
})