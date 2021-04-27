let gitSyncConfigTemplate = Vue.extend({
    template: `<div>
    <div class='row' style='line-height:60px;'>
        <div class='col-lg-2 text-left h4' style='line-height:60px;'>脚本管理</div>
        <div class='col-lg-10 h5' style='padding-right:2em;'></div>
    </div>
    <div class="row">
        <div class="col-lg-12 text-danger">
            请谨慎使用服务器上的命令,越灵活导致越不安全,更不要把账号给其他人使用.<br/>
            如果要删除不必要的文件,请使用文件管理功能,而不是使用rm命令.<br/>
            由自己造成的误操作导致的损失由自己本人负责.<br/>
            建议所有命令使用相对路径, 这样就不会导致无意中修改了webapp目录之外的文件.
        </div>
    </div>
    <hr/>
    <div>
        <form id="syncform" action="/home/GitSyncConfig" method="POST" class="form-horizontal" role="form">
            <div class="form-group">
                <label class="control-lable col-sm-2">站点</label>
                <div class="col-sm-10">
                    <input name="domain" class="form-control" type="text" v-model="domain" readonly>
                </div>
            </div>
            <div class="form-group">
                <label class="control-lable col-lg-2">初始化脚本</label>
                <div class="col-lg-10" title="先将编译的程序push到git仓库中,仅在保存时执行1次,不需要每次都上传最新编译的程序">
                    <input name="InitShell" class="form-control" type="text" v-model="init_shell" placeholder="git clone https://github.com/p4gefau1t/trojan-go.git">
                </div>
            </div>
            <div class="form-group">
                <label class="control-lable col-lg-2">定时同步脚本</label>
                <div class="col-lg-10">
                    <input name="SyncShell" class="form-control" type="text" v-model="sync_shell" placeholder="cd ./trojan-go; git pull">
                </div>
            </div>
            <div class="form-group">
                <label class="control-lable col-lg-2">同步间隔(秒)</label>
                <div class="col-lg-10">
                    <input name="interval" class="form-control" type="number" v-model="interval" min="10" max="600">
                </div>
            </div>
            <div class="form-group">
                <label class="control-lable col-lg-2">启动脚本</label>
                <div class="col-lg-10" >
                    <div><input name="StartShell" class="form-control" type="text" v-model="start_shell" placeholder="pkill trojan-go; cd ./trojan-go; chmod 777 ./trojan-go; ./trojan-go -c ./config.json"></div>
                    <div class="text-info">windows下不支持用cmd执行脚本启动程序,否则无法结束进程,建议使用指定相对路径的方式启动(不支持命令行参数,可以考虑使用配置文件), 例如: \\trojan-go\\trojan-go.exe</div>
                </div>
            </div>
            <div class="form-group">
                <label class="control-lable col-lg-2">动态密码</label>
                <div class="col-lg-10">
                    <input name="VerificationCode" class="form-control" type="text" v-model="verification_code" placeholder="Google令牌动态数字密码,此功能敬请期待">
                </div>
            </div>

            <div class="form-group">
                <div class="col-sm-10 col-sm-offset-2">
                    <button type="button" class="btn btn-primary" @click='submit($event)'>保存</button>&emsp;
                    <button type="button" class="btn btn-primary" @click='goback'>返回</button>
                </div>
            </div>
        </form>
    </div>
        </div>`,
    data: function () {
        return {
            domain: "",
            init_shell: "",
            sync_shell: "",
            interval: 60,
            start_shell: "",
            verification_code: ""
        }
    },
    methods: {
        submit: function (event) {
            $(event.target).parents("form").submit();
        },
        goback: function () {
            this.$router.go(-1);
        }
    },
    mounted: function () {
        var _this = this;
        this.domain = this.$route.query.domain;
        $.ajax({
            type: "get",
            url: "/home/GitSyncConfig",
            datatype: 'json',
            success: function (resp) {
                if (resp.code == 200 && resp.state == true && resp.data != null && resp.data != "null" && resp.data.length > 0) {
                    var configMap = JSON.parse(resp.data);
                    var config = configMap[_this.domain];
                    _this.init_shell = config.init_shell;
                    _this.sync_shell = config.sync_shell;
                    _this.interval = config.interval;
                    _this.start_shell = config.start_shell;
                }
            }
        });

        $("#syncform").ajaxForm({
            dataType: "json",
            success: function (data) {
                if (data.state == true && data.code == 200) {
                    _this.$router.go(-1);
                } else {
                    // 应该使用sweetalert
                    console.log(data);
                    alert("修改失败");
                }
            }
        });
    }
})