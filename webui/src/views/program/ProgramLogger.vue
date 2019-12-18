<template>
    <div class="h-100">
        <VTitle title-class="fa-file-text" :title="'日志查看【' + name + '】'"/>
        <div class="logger p-3">
            <div v-for="(line) in lines">{{line}}</div>
        </div>
    </div>
</template>
<script>
    import VTitle from "../../plugins/vTitle";

    export default {
        name: "ProgramLogger",
        components: {VTitle},
        data: () => ({
            node: "",
            name: "",
            lines: [],
        }),
        mounted() {
            this.node = this.$route.query.node;
            this.name = this.$route.query.name;
            this.init();
        },
        destroyed() {
            this.socket.close();
        },
        methods: {
            init() {
                if (typeof (WebSocket) === "undefined") {
                    alert("您的浏览器不支持socket")
                } else {
                    // 实例化socket
                    this.socket = new WebSocket("ws://127.0.0.1:5984/admin/program/logs?name=" + this.name + "&node=" + this.node);
                    // 监听socket连接
                    this.socket.onopen = this.open;
                    // 监听socket消息
                    this.socket.onmessage = this.getMessage;
                }
            },
            getMessage(msg) {
                this.lines.push(msg.data);
            }
        }
    }
</script>
