<template>
    <div class="h-100">
        <v-title title="系统管理" title-class="icon-settings"/>
        <div class="animated fadeIn">
            <b-tabs class="h-100" content-class="h-100">
                <b-tab active>
                    <template slot="title">
                        <i class="icons cui-envelope-closed"/>&nbsp;邮箱配置
                    </template>
                    <email/>
                </b-tab>
                <b-tab>
                    <template slot="title"><i class="icon-bell"/>&nbsp;WebHook</template>
                    <webhook/>
                </b-tab>
                <b-tab>
                    <template slot="title">
                        <i class="icons cui-comment-square"/>&nbsp;gNotify（开发中）
                    </template>
                    <h4>
                        正在开发中，敬请期待......
                    </h4>
                </b-tab>
                <b-tab>
                    <template slot="title">
                        <i class="icons cui-comment-square"/>&nbsp;推送配置说明
                    </template>
                    <h4 class="p-5">
                        <a href="https://gitee.com/ihaiker/sudis/blob/master/docs/NOTIFY.md" target="_blank">点击查看详细配置说明</a>
                    </h4>
                </b-tab>
            </b-tabs>
        </div>
    </div>
</template>

<script>
    import vTitle from "../../plugins/vTitle";
    import Email from "./email";
    import Webhook from "./webhook";
    export default {
        name: "Tags",
        components: {Webhook, Email, vTitle}
    }
</script>
