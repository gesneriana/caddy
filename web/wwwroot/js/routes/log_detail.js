let logDetailTemplate = Vue.extend({
    template: `<div>
    <div class='row' style='line-height:60px;'>
        <div class='col-lg-6 text-left h4' style='line-height:60px;'><span style='color:green;' v-bind:title='this.$route.params.TraceID'>{{this.projectName}}</span>日志详情</div>
        <div class='col-lg-6 text-right' style='padding-right:2em;'></div>
    </div>
    <div class='container'>
        <div>
            <div class='row'>
                <div class='col-lg-3'>项目ID</div>
                <div class='col-lg-9'>{{this.logTraceListData.ProjectID}}</div>
            </div>
            <div class='row'>
                <div class='col-lg-3'>Trace ID</div>
                <div class='col-lg-9'>{{this.logTraceListData.TraceID}}</div>
            </div>
            <div class='row'>
                <div class='col-lg-3'>程序集名称</div>
                <div class='col-lg-9'>{{this.logTraceListData.AssemblyName}}</div>
            </div>
            <div class='row'>
                <div class='col-lg-3'>进程名称</div>
                <div class='col-lg-9'>{{this.logTraceListData.ProcessName}}</div>
            </div>
            <div class='row'>
                <div class='col-lg-3'>进程id</div>
                <div class='col-lg-9'>{{this.logTraceListData.ProcessID}}</div>
            </div>
            <div class='row'>
                <div class='col-lg-3'>Http码</div>
                <div class='col-lg-9'>{{this.logTraceListData.HTTPStatusCode}}</div>
            </div>
            <div class='row'>
                <div class='col-lg-3'>日志级别</div>
                <div class='col-lg-9'>{{this.logTraceListData.LogLevel}}</div>
            </div>
            <div class='row'>
                <div class='col-lg-3'>机器ip</div>
                <div class='col-lg-9'>{{this.logTraceListData.ClientIP}}</div>
            </div>
            <div class='row'>
                <div class='col-lg-3'>标题</div>
                <div class='col-lg-9'>{{this.logTraceListData.LogTitle}}</div>
            </div>
            <div class='row'>
                <div class='col-lg-3'>内容</div>
                <div class='col-lg-9'>{{this.logTraceListData.LogContent}}</div>
            </div>
            <div class='row'>
                <div class='col-lg-3'>错误类型</div>
                <div class='col-lg-9'>{{this.logTraceListData.ErrorType}}</div>
            </div>
            <div class='row'>
                <div class='col-lg-3'>堆栈</div>
                <div class='col-lg-9'>{{this.logTraceListData.StackTrace}}</div>
            </div>
            <div class='row'>
                <div class='col-lg-3'>方法名</div>
                <div class='col-lg-9'>{{this.logTraceListData.FunctionName}}</div>
            </div>
            <div class='row'>
                <div class='col-lg-3'>包名</div>
                <div class='col-lg-9'>{{this.logTraceListData.PackageName}}</div>
            </div>
            <div class='row'>
                <div class='col-lg-3'>备注</div>
                <div class='col-lg-9'>{{this.logTraceListData.Remark}}</div>
            </div>
            <div class='row'>
                <div class='col-lg-3'>时间</div>
                <div class='col-lg-9'>{{this.logTraceListData.CreateTime}}</div>
            </div>
            <div style='height:36px;'></div>
        </div>
        <div v-for="(trace,index) in this.logTraceListData.LogTraceList">
            <div class='row'>
                <div class='col-lg-3'>Rpc ID</div>
                <div class='col-lg-9'>{{trace.RPCTraceID}}</div>
            </div>
            <div class='row'>
                <div class='col-lg-3'>程序集名称</div>
                <div class='col-lg-9'>{{trace.AssemblyName}}</div>
            </div>
            <div class='row'>
                <div class='col-lg-3'>进程名称</div>
                <div class='col-lg-9'>{{trace.ProcessName}}</div>
            </div>
            <div class='row'>
                <div class='col-lg-3'>操作系统</div>
                <div class='col-lg-9'>{{trace.OSName}}</div>
            </div>
            <div class='row'>
                <div class='col-lg-3'>进程id</div>
                <div class='col-lg-9'>{{trace.ProcessID}}</div>
            </div>
            <div class='row'>
                <div class='col-lg-3'>机器ip</div>
                <div class='col-lg-9'>{{trace.ClientIP}}</div>
            </div>
            <div class='row'>
                <div class='col-lg-3'>标题</div>
                <div class='col-lg-9'>{{trace.LogTitle}}</div>
            </div>
            <div class='row'>
                <div class='col-lg-3'>内容</div>
                <div class='col-lg-9'>{{trace.LogContent}}</div>
            </div>
            <div class='row'>
                <div class='col-lg-3'>错误类型</div>
                <div class='col-lg-9'>{{trace.ErrorType}}</div>
            </div>
            <div class='row'>
                <div class='col-lg-3'>堆栈</div>
                <div class='col-lg-9'>{{trace.StackTrace}}</div>
            </div>
            <div class='row'>
                <div class='col-lg-3'>方法名</div>
                <div class='col-lg-9'>{{trace.FunctionName}}</div>
            </div>
            <div class='row'>
                <div class='col-lg-3'>包名</div>
                <div class='col-lg-9'>{{trace.PackageName}}</div>
            </div>
            <div class='row'>
                <div class='col-lg-3'>备注</div>
                <div class='col-lg-9'>{{trace.Remark}}</div>
            </div>
            <div class='row'>
                <div class='col-lg-3'>时间</div>
                <div class='col-lg-9'>{{trace.CreateTime}}</div>
            </div>
            </br>
        </div>
    </div>
    </div>`,
    data: function () {
        return {
            projectName: "项目名称",
            logTraceListData: {}
        }
    },
    methods: {
    },
    mounted: function () {
        // console.log(this.$route.params.TraceID)
        this.projectName = $.cookie('CurrentProjectName');
        window.projectThis = this;
        $.ajax({
            type: "get",
            url: "/home/logtrace?traceid=" + this.$route.params.TraceID,
            datatype: 'json',
            success: function (resp) {
                if (resp.state == true) {
                    window.projectThis.logTraceListData = resp.data;
                }
            }
        });
    }
})