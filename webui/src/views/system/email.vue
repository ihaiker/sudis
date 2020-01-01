<template>
    <div>
        <div class="alert alert-info">邮件服务器配置</div>
        <div class="row ">
            <div class="col-4 form-group">
                <div class="input-group">
                    <div class="input-group-prepend">
                        <span class="input-group-text">
                            <i class="fa fa-user"/> SMTP服务地址：
                        </span>
                    </div>
                    <input class="form-control" v-model="address" type="text" placeholder="SMTP地址"/>
                </div>
            </div>
            <div class="col-4 form-group">
                <div class="input-group">
                    <div class="input-group-prepend">
                        <span class="input-group-text">
                            <i class="fa fa-envelope"/>
                        </span>
                    </div>
                    <input class="form-control" v-model="port" type="number" placeholder="SMTP端口号"/>
                </div>
            </div>
            <div class="col-4">
                <button class="btn btn-primary w-25" @click="setConfig">设置</button>
            </div>
            <div class="col"></div>
        </div>

        <div class="row">
            <div class="col form-group">
                <div class="input-group">
                    <div class="input-group-prepend">
                        <span class="input-group-text">
                            <i class="fa fa-user"/>&nbsp;用户名
                        </span>
                    </div>
                    <input class="form-control" v-model="name" type="text" placeholder="Username">
                </div>
            </div>
            <div class="col form-group">
                <div class="input-group">
                    <div class="input-group-prepend">
                        <span class="input-group-text">
                            <i class="fa fa-asterisk"/>&nbsp;密码
                        </span>
                    </div>
                    <input class="form-control" v-model="passwd" type="text" placeholder="Password">
                </div>
            </div>
            <div class="col">
                <button class="btn btn-tumblr w-25" @click="testConfig">测试</button>
            </div>
        </div>
        <div class="row">
            <div class="col form-group">
                <div class="input-group">
                    <div class="input-group-prepend">
                        <span class="input-group-text">
                            <i class="fa fa-user"/>&nbsp;接收用户
                        </span>
                    </div>
                    <input class="form-control" v-model="to" type="text" placeholder="接收用户">
                </div>
            </div>
        </div>

        <div class="alert alert-info mt-2">通知模板</div>
        <attrs @change="content = content + $event"/>
        <textarea class="form-control mt-3" style="min-height: 300px;" v-model="content"/>
    </div>
</template>

<script>
    import Attrs from "./attrs";

    export default {
        name: "email",
        components: {Attrs},
        data: () => ({
            address: "", port: 465,
            name: "", passwd: "", to: "",
            content: `节点：{{.Node}}，程序：{{.Name}}，状态更改：{{.State}}`,
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
                    self.$toast.error(e.message);
                })
            },
            testConfig() {
                this.execConfig("/admin/notify/test");
            },
            setConfig(){
                this.execConfig("/admin/notify");
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
                    self.$toast.error(e.message);
                });
            }
        }
    }
</script>
