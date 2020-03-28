<template>
    <div>
        <div class="alert alert-info"><i class="icon-settings"/>&nbsp;地址设置：WebHook通知。Webhook将以POST形式推送模板JSON内容</div>
        <div class="input-group">
            <div class="input-group-prepend">
                <span class="input-group-text">
                    <i class="fa fa-user"/>&nbsp;WebHook地址：
                </span>
            </div>
            <input class="form-control" v-model="address" type="text" placeholder="WebHook地址"/>
            <div class="input-group-append">
                <button class="btn btn-block btn-primary" style="width: 120px" @click="setConfig">设置</button>
            </div>
        </div>
        <div class="input-group mt-2">
            <div class="input-group-prepend">
                <span class="input-group-text">
                    <i class="fa fa-asterisk"/>&nbsp;Access Token：
                </span>
            </div>
            <input class="form-control" v-model="token" type="text" placeholder="SecurityToken"/>
            <div class="input-group-append">
                <button class="btn btn-block btn-vine" style="width: 120px" @click="testConfig">测试</button>
            </div>
            <div class="input-group-append">
                <button class="btn btn-outline-danger" style="width: 120px" @click="clearConfig">清除</button>
            </div>
        </div>
        <div class="alert alert-info mt-2">
            通知模板：<span v-if="address">{{address}}?access_token={{token}}</span>
        </div>
        <textarea class="form-control mt-3" style="min-height: 350px;" v-model="content" placeholder="使用程序默认JSON推送"/>
    </div>
</template>

<script>
    export default {
        name: "webhook",
        data: () => ({
            address: "", token: "", content: ``
        }),
        mounted() {
            this.getConfig();
        },
        methods: {
            testConfig() {
                this.execConfig("/admin/notify/test")
            },
            setConfig() {
                this.execConfig("/admin/notify")
            },
            getConfig() {
                let self = this;
                self.$axios.get("/admin/notify/webhook").then(res => {
                    let config = JSON.parse(res.config);
                    self.address = config.address;
                    self.token = config.token;
                    self.content = config.content;
                }).catch(e => {
                    self.$toast.error("WebHock" + e.message);
                })
            },
            clearConfig(){
                let self = this;
                self.$axios.delete("/admin/notify/webhook").then(res => {
                    self.$toast.success('清除成功！');
                    self.address = "";
                    self.token = "";
                    self.content = "";
                }).catch(e => {
                    self.$toast.error("email" + e.message);
                });
            },
            execConfig(uri) {
                let self = this;
                let config = {address: self.address, token: self.token, content: self.content};
                self.$axios.post(uri, {name: "webhook", config: JSON.stringify(config)}).then(res => {
                    self.$toast.success('成功！');
                }).catch(e => {
                    self.$toast.error("WebHock" + e.message);
                });
            },
        }
    }
</script>
