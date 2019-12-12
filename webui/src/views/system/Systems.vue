<template>
    <div class="h-100">
        <v-title title="系统管理" title-class="icon-settings"/>
        <div class="alert alert-danger">
            功能暂未实现，一人利用工作之余时间开发，时间紧迫，如果您对本程序的开发有兴趣可联系。 <a href="mailto:ni@renzhen.la">Haiker (ni@renzhen.la)</a>
        </div>
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

<style scoped>

</style>
