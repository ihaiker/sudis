<template>
    <div class="btn-group btn-group-sm" role="group" aria-label="Small button group">
        <button class="btn btn-default" v-for="(name,attr) in attrs"
                @click="setAttr(attr)">{{attr}}：{{name}}
        </button>
    </div>
</template>

<script>
    export default {
        name: "attrs",
        data: () => ({
            attrs: {
                Id: "程序ID", Node: "节点", Name: "程序名称", WorkDir: "工作目录", User: "运行用户",
                Envs: "环境变量", State: "状态",
                "Start.Command": "启动命令", "Start.Args": "启动命令参数",
                "Stop.Command": "停止命令", "Stop.Args": "停止命令参数",
            },
        }),
        methods: {
            setAttr(attr) {
                this.$emit("change", "{{." + attr + "}}");
            }
        }
    }
</script>
