let jsonConfigTemplate = Vue.extend({
    template:
    `<div>
        <div class='row' style='line-height:60px;'>
            <div class='col-lg-3 text-left h4' style='line-height:60px;'>
                <div>Caddy配置文件<a href="https://caddyserver.com/docs/caddyfile/concepts" target="_blank">帮助文档</a></div>
            </div>
            <div class='col-lg-10 text-right' style='padding-right:2em;'></div>
        </div>
    
        <div style='line-height:60px;'>
            <ul class="nav nav-tabs row" role="tablist">
                <li role="presentation" class="active col-lg-1" >
                    <a href="#menu1" role="tab" data-toggle="tab">Caddy File</a>
                </li>
                <li role="presentation" class="col-lg-1" >
                    <a href="#menu2" v-on:click="SetCaddyJsonRaw" role="tab" data-toggle="tab">Raw</a>
                </li>
                <li role="presentation" class="col-lg-1">
                    <a href="#menu3" v-on:click="SetEditor" role="tab" data-toggle="tab">Editor</a>
                </li>
                <li role="presentation" class="col-lg-1">
                    <a href="#menu4" role="tab" data-toggle="tab">模板</a>
                </li>
            </ul>
        </div>
    
        <div class="tab-content">
            <div id="menu1" role="tabpanel" class="tab-pane active">
                <div class="row">
                    <div class='col-lg-11'>
                        <textarea id="txtCaddy" v-model="caddyConfig" style="width: 100%; min-height: 600px"></textarea>
                    </div>
                </div>
                <div class="row">
                    <div class='col-lg-11'>
                        <p>为什么采用caddyfile进行配置, 而不是全部使用json api, 请参考<a href="https://dengxiaolong.com/caddy/v2/zh/getting-started.html" target="_blank">对比说明</a></p>
                        <p>caddyfile配置永久生效,JSON API只在caddy服务器运行时生效,重启后只会读取caddyfile的配置文件</p>
                        <p>caddyfile已经能满足基本的功能,并且使用方便, 虽然JSON API包含所有的功能和模块,但是过于复杂</p>
                    </div>
                </div>
                <div class='row' style='line-height:60px;'>
                    <div class='col-lg-10'></div>
                    <div class='col-lg-1 text-right'>
                        <button class='btn btn-primary' v-on:click='SaveCaddyFileConfig' >保存</button>
                    </div>
                </div>
            </div>
    
            <div id="menu2" role="tabpanel" class="tab-pane">
                <textarea id="txtRaw" v-model="jsonConfig" style="width: 90%; min-height: 800px"></textarea>
            </div>
            <div id="menu3" role="tabpanel" class="tab-pane">
                <div class='row'>
                    <div id="jsoneditor" class='col-lg-11' style="min-height: 800px;"></div>
                </div>
                <div class='row' style='line-height:60px;'>
                    <div class='col-lg-10'></div>
                    <div class='col-lg-1 text-right'>
                        <button class='btn btn-primary' v-on:click='SaveCaddyJsonConfig' >保存</button>
                    </div>
                </div>
            </div>
            <div id="menu4" role="tabpanel" class="tab-pane">
                <textarea id="txtDemo" v-model="demo" style="width: 90%; min-height: 800px" readonly></textarea>
            </div>
        </div>
    
    </div>`,
    data: function () {
        return {
            editor: {},
            caddyConfig: "",
            jsonConfig: "",
            demo: "",
        }
    },
    methods: {
        SetCaddyJsonRaw: function () {
            // get json
            const updatedJson = this.editor.get();
            if ($.isEmptyObject(updatedJson)) {
                return;
            }
            this.jsonConfig = JSON.stringify(updatedJson, null, '\t');
        },
        SetEditor: function () {
            if (this.jsonConfig != null && this.jsonConfig.length > 0) {
                this.editor.set(JSON.parse(this.jsonConfig))
            }
        },
        SaveCaddyJsonConfig: function () {
            // get json
            const updatedJson = this.editor.get();
            if ($.isEmptyObject(updatedJson)) {
                console.log("请输入有效的json配置文件");
                return;
            }

            console.log("保存配置文件");
            console.log(updatedJson);
            $.ajax({
                type: "post",
                url: "/caddy/json_config",
                contentType: 'application/json',  //指定格式为json格式
                datatype: 'json',
                data: JSON.stringify(updatedJson),
                success: function (resp) {
                    console.log(resp);
                    if (resp.code == 200 && resp.data != null && resp.data != "null") {
                        console.log(resp.data)
                    }
                }
            });
        },
        SaveCaddyFileConfig: function () {
            // get json
            var config = this.caddyConfig;
            if (config == null || config.length == 0) {
                console.log("请输入有效的caddy配置文件");
                return;
            }

            console.log("保存配置文件");
            console.log(config);
            $.ajax({
                type: "post",
                url: "/caddy/caddy_config",
                datatype: 'json',
                data: "Caddy=" + config,
                success: function (resp) {
                    console.log(resp);
                    if (resp.code == 200 && resp.data != null && resp.data != "null") {
                        console.log(resp.data)
                    }
                }
            });
        }
    },
    mounted: function () {
        // create the editor
        const container = document.getElementById("jsoneditor")
        const options = {}
        this.editor = new JSONEditor(container, options)
        var jsonEdtior = this.editor;
        var _this = this;

        $.ajax({
            type: "get",
            url: "/caddy/json_config",
            datatype: 'json',
            success: function (resp) {
                // console.log(resp);
                if (resp.code == 200 && resp.data != null && resp.data != "null" && resp.data != "null\n") {
                    _this.jsonConfig = resp.data;
                    jsonEdtior.set(JSON.parse(_this.jsonConfig));
                }
            }
        });

        $.ajax({
            type: "get",
            url: "/caddy/caddy_config",
            datatype: 'json',
            success: function (resp) {
                if (resp.code == 200 && resp.data != null && resp.data != "null" && resp.data != "null\n") {
                    _this.caddyConfig = resp.data;
                }
            }
        });

        this.demo = `caddy.kizuna.top {
            reverse_proxy 127.0.0.1:2020
        }
        
        # jp.kizuna.top {
        #    reverse_proxy 127.0.0.1:2020 127.0.0.1:2021 {
        #        lb_policy first
        #    }
        # }`
    }
})