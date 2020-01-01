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
            attrs: {Node: "节点", Name: "程序名称", State: "状态",},
        }),
        methods: {
            setAttr(attr) {
                this.$emit("change", "{{." + attr + "}}");
            }
        }
    }
</script>
