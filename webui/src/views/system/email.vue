<template>
    <div>
        <div class="alert alert-info">邮件服务器配置</div>

        <div class="input-group">
            <div class="input-group-prepend">
                <span class="input-group-text">
                    <i class="fa fa-user"/> SMTP服务地址：
                </span>
            </div>
            <input class="form-control" v-model="address" type="text" placeholder="SMTP地址"/>
            <div class="input-group-prepend">
                <span class="input-group-text">
                    <i class="fa fa-envelope"/>
                </span>
            </div>
            <input class="form-control" v-model="port" type="number" placeholder="SMTP端口号"/>
            <div class="input-group-append">
                <button class="btn btn-primary" @click="setConfig"> 设 置</button>
            </div>
            <div class="input-group-append">
                <button class="btn btn-dark" @click="clearConfig">清除设置</button>
            </div>
        </div>

        <div class="input-group mt-2">
            <div class="input-group-prepend">
                <span class="input-group-text">
                    <i class="fa fa-user"/>&nbsp;用户名
                </span>
            </div>
            <input class="form-control" v-model="name" type="text" placeholder="Username">
            <div class="input-group-prepend">
                <span class="input-group-text">
                    <i class="fa fa-asterisk"/>&nbsp;密码
                </span>
            </div>
            <input class="form-control" v-model="passwd" type="text" placeholder="Password">
            <div class="input-group-prepend">
                <span class="input-group-text">
                    <i class="fa fa-user"/>&nbsp;接收用户
                </span>
            </div>
            <input class="form-control" v-model="to" type="text" placeholder="接收用户">
            <div class="input-group-append">
                <button class="btn btn-tumblr w-100" @click="testConfig">测试</button>
            </div>
        </div>

        <div class="alert alert-info mt-2">通知模板</div>
        <textarea class="form-control mt-3" style="min-height: 300px;" v-model="content"/>
    </div>
</template>

<script>
    export default {
        name: "email",
        data: () => ({
            address: "", port: 465,
            name: "", passwd: "", to: "",
            content: `
{{ if eq .Type "process" }}
    <h1>程序通知：</h1>
    <table>
        <tr>
            <th>程序节点</th><td>{{.Process.Node}}</td>
            <th>名称</th><td>{{.Process.Name}}</td>
        </tr>
        <tr>
            <th>From:</th><td>{{.FromStatus}}</td>
            <th>To:</th><td>{{.ToStatus}}</td>
        </tr>
    </table>
{{else}}
    <h1>节点通知：</h1>
    {{.Node}} 状态变更为：{{.Status}}
{{end}}
`,
        }),
        mounted() {
            this.getConfig();
        },
        methods: {
            getConfig() {
                let self = this;
                self.$axios.get("/admin/notify/email").then(res => {
                    let config = JSON.parse(res.config);
                    self.address = config.address;
                    self.port = config.port;
                    self.name = config.name;
                    self.passwd = config.passwd;
                    self.content = config.content;
                    self.to = config.to;
                }).catch(e => {
                    self.$toast.error("email" + e.message);
                })
            },
            testConfig() {
                this.execConfig("/admin/notify/test");
            },
            setConfig() {
                this.execConfig("/admin/notify");
            },
            clearConfig() {
                let self = this;
                self.$axios.delete("/admin/notify/email").then(res => {
                    self.$toast.success('清除成功！');
                    self.address = "";
                    self.name = "";
                    self.passwd = "";
                    self.content = "";
                    self.to = "";
                }).catch(e => {
                    self.$toast.error("email" + e.message);
                });
            },
            execConfig(uri) {
                let self = this;
                let config = {
                    address: self.address, port: parseInt(self.port),
                    name: self.name, passwd: self.passwd, content: self.content,
                    to: self.to,
                };
                self.$axios.post(uri, {name: "email", config: JSON.stringify(config)}).then(res => {
                    self.$toast.success('成功！');
                }).catch(e => {
                    self.$toast.error("email" + e.message);
                });
            }
        }
    }
</script>
