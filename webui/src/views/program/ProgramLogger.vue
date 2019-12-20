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
            getDomain() {
                //如果配置里，就是用配置地址
                let configDomain = process.env.VUE_APP_WS;
                if (configDomain !== undefined && configDomain !== '') {
                    return configDomain;
                }
                //如果未配置，测试网使用默认本机地址，其他根据域名解析出ws地址
                if (process.env.NODE_ENV === "development") {
                    return "ws://127.0.0.1:5984"
                } else {
                    let domain = window.location.href.substr(0, window.location.href.indexOf("/", 8));
                    domain = domain.replace("https://", "wss://");
                    domain = domain.replace("http://", "ws://");
                    return domain;
                }
            },

            init() {
                if (typeof (WebSocket) === "undefined") {
                    alert("您的浏览器不支持socket")
                } else {
                    let self = this;
                    this.socket = new WebSocket(self.getDomain() + "/admin/program/logs");
                    this.socket.onopen = () => {
                        self.socket.send(JSON.stringify({
                            name: this.name,
                            node: this.node,
                            user: localStorage.getItem("x-user"),
                            ticket: localStorage.getItem("x-ticket")
                        }));
                    };
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
