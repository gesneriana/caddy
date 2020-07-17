let certlistTemplate = Vue.extend({
    template: `<div>
    <div class='row' style='line-height:60px;'>
        <div class='col-lg-1 text-left h4' style='line-height:60px;'>证书管理</div>
        <div class='col-lg-9 text-left h6' style='line-height:60px;'></div>
        <div class='col-lg-1 text-right'></div>
    </div>
    <div class="row">
        <div class='col-lg-11'>
        <table class="table table-bordered table-striped" id="sitelist-table">
            <tr>
                <th style="" data-field="handle"><div class="th-inner ">域名</div><div class="fht-cell"></div></th>
                <th style="" data-field="upsteams"><div class="th-inner ">目录</div><div class="fht-cell"></div></th>
                <th style="" data-field="upsteams"><div class="th-inner ">文件名</div><div class="fht-cell"></div></th>
                <th style="" data-field="match-value"><div class="th-inner ">修改时间</div><div class="fht-cell"></div></th>
                <th style="text-align: center; " data-field="index"><div class="th-inner ">编辑</div><div class="fht-cell"></div></th>
            </tr>
            <tr v-for="(item, index) in certData.CertList">
                <td>
                    <span>{{item.Domain}}</span>
                </td>
                <td>
                    <span>{{item.CertDir}}</span>
                </td>
                <td>
                    <span>{{item.CertDir}}/{{item.Domain}}.crt</span><br/>
                    <span>{{item.CertDir}}/{{item.Domain}}.key</span>
                </td>
                <td>
                    <span>{{item.LastModified}}</span>
                </td>
                <td class="text-center">
                    <button class='btn btn-primary' @click="exportCert(index)">导出</button>
                </td>
            </tr>
        </table>
        </div>
    </div>
    </div>`,
    data: function () {
        return {
            certData: {}
        }
    },
    methods: {
        exportCert: function (index) {
            console.log("下载证书文件" + index);
        },
    },
    mounted: function () {
        var _this = this;
        $.ajax({
            type: "get",
            url: "/caddy/certlist",
            datatype: 'json',
            success: function (resp) {
                if (resp.code == 200 && resp.state == true) {
                    _this.certData = JSON.parse(resp.data);
                    // console.log(_this.caddyRoutes);
                }
            }
        });
    }
})