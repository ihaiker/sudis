<template>
    <div class="row">
        <div class="col-auto">
            <div class="form-group">
                <div class="input-group input-group-sm">
                    <div class="input-group-prepend">
                        <span class="input-group-text">进程名</span>
                    </div>
                    <input class="form-control" @keyup.enter="onSearch" v-model="form.name" type="text" name="program" placeholder="进程名称"/>
                </div>
            </div>
        </div>
        <div class="col-auto pl-0">
            <div class="form-group">
                <div class="input-group input-group-sm">
                    <span class="input-group-prepend">
                        <span class="input-group-text">所属节点</span>
                    </span>
                    <select class="form-control" v-model="form.node">
                        <option value="">所有节点　　　　　　　　　</option>
                        <option v-for="node in nodes" :value="node.key">{{node.tag === "" ? node.key : node.tag }}</option>
                    </select>
                </div>
            </div>
        </div>
        <div class="col-auto pl-0">
            <div class="form-group">
                <div class="input-group input-group-sm">
                    <div class="input-group-prepend">
                        <span class="input-group-text">含有标签</span>
                    </div>
                    <select class="form-control" v-model="form.tag">
                        <option value="">所有标签</option>
                        <option v-for="tag in tags" :value="tag.name">{{tag.name}}</option>
                    </select>
                </div>
            </div>
        </div>
        <div class="col-auto pl-0">
            <div class="form-group">
                <div class="input-group input-group-sm">
                    <div class="input-group-prepend">
                        <span class="input-group-text">当前状态</span>
                    </div>
                    <select class="form-control" v-model="form.status">
                        <option value="">所有状态</option>
                        <option value="ready">准备</option>
                        <option value="running">运行</option>
                        <option value="stoped">停止</option>
                        <option value="fail">异常</option>
                    </select>
                </div>
            </div>
        </div>
        <div class="col-auto pl-0">
            <button class="btn btn-sm btn-default" @click="onSearch">
                <i class="fa fa-search"/> 搜索
            </button>
            <button class="btn btn-sm btn-default ml-2" @click="onReset">
                <i class="fa fa-close"/> 重置搜索
            </button>
            <slot/>
        </div>


    </div>
</template>

<script>
    import Multiselect from 'vue-multiselect'
    import 'vue-multiselect/dist/vue-multiselect.min.css'

    export default {
        name: "search",
        components: {Multiselect},
        props: {
            nodes: {
                type: Array,
                default: [],
            },
            tags: {
                type: Array,
                default: [],
            },
        },
        data: () => ({
            form: {name: "", node: "", tag: "", status: ""},
            options: [
                {name: "所有状态", value: ""},
                {name: "准备", value: "ready"},
                {name: "运行", value: "running"},
                {name: "停止", value: "stoped"},
                {name: "异常", value: "fail"},
            ],
        }),
        mounted() {
            let tag = this.$route.params.tag;
            if (tag) {
                this.form.tag = tag;
            }
            this.onSearch()
        },
        methods: {
            onSearch() {
                this.$emit("search", this.form);
            },
            onReset() {
                this.form = {name: "", node: "", tag: "", status: ""};
                this.onSearch();
            }
        },
    }
</script>
