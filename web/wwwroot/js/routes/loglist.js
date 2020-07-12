let loglistTemplate = Vue.extend({
    template: `<div>
    <div class='row' style='line-height:60px;'>
        <div class='col-lg-6 text-left h4' style='line-height:60px;'><span style='color:green;' v-bind:title='this.$route.params.ProjectID'>{{this.projectName}}</span>日志列表</div>
        <div class='col-lg-6 text-right' style='padding-right:2em;'></div>
    </div>
    <div class='container'>
        <div class='hide' v-for="(log,index) in this.logdata">
            <div class='row'>
                <div class='col-lg-3'>项目ID</div>
                <div class='col-lg-9'>{{log.ProjectID}}</div>
            </div>
            <div class='row'>
                <div class='col-lg-3'>日志ID</div>
                <div class='col-lg-9'>{{log.TraceID}}</div>
            </div>
            <div class='row'>
                <div class='col-lg-3'>AssemblyName</div>
                <div class='col-lg-9'>{{log.AssemblyName}}</div>
            </div>
            <div class='row'>
                <div class='col-lg-3'>ProcessName</div>
                <div class='col-lg-9'>{{log.ProcessName}}</div>
            </div>
        </div>
        <table class="table table-bordered table-striped" data-show-refresh="true" id="logdata-table"></table>
    </div>
    </div>`,
    data: function () {
        return {
            projectName: "项目名称",
            logdata: [],
            columns: [
                { title: '日志ID', field: 'TraceID' },
                { title: '日志等级', field: 'LogLevel' },
                { title: '标题', field: 'LogTitle' },
                { title: '创建时间', field: 'CreateTime' },
                {
                    title: '查看详情',
                    field: 'ProjectID',
                    align: 'center',
                    formatter: function (value, row, index) {
                        var a = '<button class="btn btn-primary" onclick=\'showLogDetail("' + row.TraceID + '")\'>查看详情</button>';
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
        showLogDetail: function (traceId) {
            // 打开详情页面
            // console.log(traceId);
            window.projectThis.$router.push({
                name: 'logdetail',
                params:
                {
                    TraceID: traceId
                }
            })
        }
    },
    mounted: function () {
        // console.log(this.$route.params.ProjectID)
        this.projectName = $.cookie('CurrentProjectName');
        window.showLogDetail = this.$options.methods.showLogDetail;
        window.projectThis = this;

        $("#logdata-table").bootstrapTable(
            {
                url: "/home/findlogdatawithpager?pid=" + this.$route.params.ProjectID,
                method: "get",
                striped: true,
                cache: true,
                pagination: true,
                sortable: true,
                sortOrder: "asc",
                sidePagination: "server",
                pageNumber: 1,
                pageList: [10, 25, 50, 100],
                search: true,
                uniqueId: "TraceID",
                queryParams: function (params) {
                    //这里的键的名字和控制器的变量名必须一致，这边改动，控制器也需要改成一样的
                    var temp = {
                        rows: params.limit,                         //页面大小
                        page: (params.offset / params.limit) + 1,   //页码
                        sort: params.sort,      //排序列名  
                        sortOrder: params.order //排位命令（desc，asc） 
                    };
                    return temp;
                },
                columns: this.columns,
                onDblClickRow: function (row, $element) {
                    console.log("双击了 " + row.TraceID);
                },
            }
        )
    }
})