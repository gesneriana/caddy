let pwdmanageTemplate = Vue.extend({
    template:
        `<div>
        <div class='row' style='line-height:60px;'>
            <div class='col-lg-3 text-left h4' style='line-height:60px;'>
                <div>Caddy密码管理</div>
            </div>
            <div class='col-lg-10 text-right' style='padding-right:2em;'></div>
        </div>
    
        <div style='line-height:60px;'>
            <ul class="nav nav-tabs row" role="tablist">
                <li role="presentation" class="active col-lg-1" >
                    <a href="#menu1" role="tab" data-toggle="tab">更改密码</a>
                </li>
                <li role="presentation" class="col-lg-1" >
                    <a href="#menu2" role="tab" data-toggle="tab" >查看Token</a>
                </li>
            </ul>
        </div >
    
        <div class="tab-content">
            <div id="menu1" role="tabpanel" class="tab-pane active">
                <div class="row">
                    <div class='col-lg-11'>
                        <div class="row">
                            <div class="col-lg-4"></div>
                            <div class="col-lg-4">
                                <div>
                                    <form id="pwdform" action="/home/changepwd" method="POST" class="form-horizontal" role="form">
                                        <div class="form-group text-center">
                                            <legend>更新用户</legend>
                                        </div>

                                        <div class="form-group">
                                            <label for="username" class="control-lable col-sm-3">用户名</label>
                                            <div class="col-sm-9">
                                                <input name="username" class="form-control" type="text" placeholder="忘记密码可以删除pwd.json配置文件"/>
                                            </div>
                                        </div>
                                        <div class="form-group">
                                            <label for="oldpassword" class="control-lable col-sm-3">原密码</label>
                                            <div class="col-sm-9">
                                                <input name="oldpassword" class="form-control" type="password" placeholder="忘记密码可以删除pwd.json配置文件"/>
                                            </div>
                                        </div>
                                        <div class="form-group">
                                            <label class="control-lable col-lg-3">新密码</label>
                                            <div class="col-lg-9">
                                                <input name="newpassword" class="form-control" type="password" />
                                            </div>
                                        </div>
    
                                        <div class="form-group">
                                            <div class="col-sm-10 col-sm-offset-3">
                                                <button type="button" v-on:click="submitPwdForm($event)" class="btn btn-primary">修改</button>
                                            </div>
                                        </div>
                                    </form>
                                </div>
    
                            </div>
                            <div class="col-lg-4"></div>
                        </div>
                    </div>
                </div>
            </div>
    
            <div id="menu2" role="tabpanel" class="tab-pane">
                <textarea id="txtRaw" v-model="token" style="width: 90%; min-height: 800px"></textarea>
            </div>
        </div>
    
    </div>`,
    data: function () {
        return {
            demo: "",
            token: "",
        }
    },
    methods: {
        submitPwdForm: function (event) {
            $(event.target).parents("form").submit();
        }
    },
    mounted: function () {
        var _this = this;

        $("#pwdform").ajaxForm({
            dataType: "json",
            success: function (data) {
                if (data.state == true && data.code == 200) {
                    // 跳转到新页面
                    location.href = "/view/login.html"
                } else {
                    // 应该使用sweetalert
                    console.log(data);
                    alert("修改失败");
                }
            }
        });

        this.demo =
            `caddy.kizuna.top {
            reverse_proxy 127.0.0.1:2020
        }
        
        # jp.kizuna.top {
        #    reverse_proxy 127.0.0.1:2020 127.0.0.1:2021 {
        #        lb_policy first
        #    }
        # }
        `
    }
})