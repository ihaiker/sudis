<template>
    <div>
        <div class="alert alert-link">WebHook通知。Webhook将以POST形式推送模板JSON内容</div>

        <div class="alert alert-info mt-2"><i class="icon-settings"/> 地址设置：</div>
        <div class="input-group">
            <div class="input-group-prepend">
                <span class="input-group-text">
                    <i class="fa fa-user"/> WebHook地址：
                </span>
            </div>
            <input class="form-control" v-model="address" type="text" placeholder="WebHook地址"/>
        </div>
        <div class="input-group mt-2">
            <div class="input-group-prepend">
                <span class="input-group-text">
                    <i class="fa fa-asterisk"/> SecurityToken：
                </span>
            </div>
            <input class="form-control" v-model="token" type="text" placeholder="SecurityToken"/>
        </div>
        <center>
            <button class="btn btn-primary w-25 mt-2">设置</button>
        </center>

        <div class="alert alert-info mt-2">通知模板：</div>
        <pre><code>
{
    name: "", node: "", daemon: "0",
    autoStart: false, ignoreStarted: true,
    startDuration: 7, startRetries: 3,
    envs: [],
    start: {
        command: "", args: [],
        health: {
            type: "", url: "", ttl: "5", securityKey: ""
        }
    },
    stop: {
        command: "", args: [],
    },
    workDir: "", user: "",
}</code></pre>
    </div>
</template>

<script>
    export default {
        name: "webhook",
        data: () => ({
            address: "", token: "",
            config: {
                mode: 'javascript',
            },
            code: `
            `
        })
    }
</script>
