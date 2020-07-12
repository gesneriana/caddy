let editProjectTemplate = Vue.extend({
    template: `<div>
    <div class='row' style='line-height:60px;'>
        <div class='col-lg-2 text-left h4' style='line-height:60px;'>编辑项目</div>
        <div class='col-lg-10 text-right' style='padding-right:2em;'></div>
    </div>
    <div>
        <form id="projectform" action="/home/editProject" method="POST" class="form-horizontal" role="form">
            <input id="ProjectID" name="ProjectID" class="form-control" type="hidden" v-model="this.$route.params.ProjectID" >
            <div class="form-group">
                <label for="projectname" class="control-lable col-sm-2">项目名称</label>
                <div class="col-sm-10">
                    <input id="projectname" name="ProjectName" class="form-control" type="text" v-model="this.$route.params.ProjectName" >
                </div>
            </div>
            <div class="form-group">
                <label class="control-lable col-lg-2">开发平台</label>
                <div class="col-lg-10">
                    <input name="ProjectPlatform" class="form-control" type="text" v-model="this.$route.params.ProjectPlatform" >
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
    methods: {
        submit: function (event) {
            $(event.target).parents("form").submit()
        }
    },
    mounted: function () {
        var _this = this;
        $("#projectform").ajaxForm({
            dataType: "json",
            success: function (data) {
                if (data.state === true) {
                    // 跳转到新页面
                    _this.$router.push({ path: '/project' })
                } else {
                    // 应该使用sweetalert
                    alert(data.Message + "," + data.Error)
                }
            }
        })
    }
})