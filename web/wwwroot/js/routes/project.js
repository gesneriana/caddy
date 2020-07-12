let projectTemplate = Vue.extend({
    template: `<div>
    <div class='row' style='line-height:60px;'>
        <div class='col-lg-2 text-left h4' style='line-height:60px;'>项目</div>
        <div class='col-lg-10 text-right' style='padding-right:2em;'><button class='btn btn-primary' @click='NewProject' >New</button></div>
    </div>
    <div>
        <table class="table table-bordered table-striped" data-show-refresh="true" id="project-table"></table>
    </div>
    </div>`,
    data: function () {
        return {
            columns: [
                { title: '项目ID', field: 'ProjectID' },
                {
                    title: '项目名称',
                    field: 'ProjectName',
                    formatter: function (value, row, index) {
                        var a = '<a href="javascript:;" onclick="showProject(\'' + row.ProjectID + '\',\'' + row.ProjectName + '\')">' + value + '</a>';
                        return a
                    }
                },
                { title: '项目开发平台', field: 'ProjectPlatform' },
                { title: '创建者', field: 'Cretor' },
                { title: '创建时间', field: 'CreateTime' },
                {
                    title: '编辑',
                    field: 'ProjectID',
                    align: 'center',
                    formatter: function (value, row, index) {
                        var a = "<button class='btn btn-primary' onclick='editProject(\"" + row.ProjectID + "\",\"" + row.ProjectName + "\",\"" + row.ProjectPlatform + "\")'>修改</button>";
                        return a
                    }
                }
            ]
        }
    },
    methods: {
        NewProject: function () {
            this.$router.push({ path: '/newproject' })
        },
        editProject: function (projectId, projectName, projectPlatform) {
            // var tds = $($("#project-table tbody tr").get(index)).children();
            // console.log(tds.innerText)
            // 打开编辑页面
            window.projectThis.$router.push({
                name: 'editproject',
                params:
                {
                    ProjectID: projectId,
                    ProjectName: projectName,
                    ProjectPlatform: projectPlatform
                }
            })
        },
        showProject: function (projectId, projectName) {
            // var tds = $($("#project-table tbody tr").get(index)).children();
            // console.log(tds.get(0).innerText)
            // 打开详情页面
            $.cookie('CurrentProjectId', projectId, { expires: 365 });
            $.cookie('CurrentProjectName', projectName, { expires: 365 });
            window.projectThis.$router.push({
                name: 'loglist',
                params:
                {
                    ProjectID: projectId
                }
            })
        }
    },
    mounted: function () {
        $("#project-table").bootstrapTable(
            {
                method: 'get',
                url: "/home/project_data",
                columns: this.columns,
                search: true,
                pagination: true
            }
        )
        window.editProject = this.$options.methods.editProject;
        window.showProject = this.$options.methods.showProject;
        window.projectThis = this;
    }
})