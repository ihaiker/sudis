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

        <div class="col-auto">
            <div class="btn-group btn-group-sm" role="group" aria-label="Small button group">
                <button class="btn" @click="onNode('')" :class="form.node === ''?'btn-primary':'btn-default'">所有节点</button>
                <button class="btn" v-for="node in nodes" :class="form.node === node.key ? 'btn-primary':'btn-default'"
                        @click="onNode(node.key)">{{node.tag === "" ? node.key : node.tag }}
                </button>
            </div>
        </div>
        <div class="col-auto">
            <div class="btn-group btn-group-sm" role="group">
                <button class="btn" :class="form.tag === ''?'btn-primary':'btn-default'" @click="onTag('')">所有标签</button>
                <button class="btn" v-for="tag in tags" :class="tag.name === form.tag? tag.class : 'btn-default' "
                        @click="onTag(tag.name)" style="min-width: 60px;">{{tag.name}}
                </button>
            </div>
        </div>

        <div class="col-auto ">
            <div class="btn-group btn-group-sm" role="group">
                <button class="btn" @click="onStatus('')" :class="form.status === '' ? 'btn-primary' : 'btn-default' ">所有状态</button>
                <button class="btn" @click="onStatus('ready')" :class="form.status === 'ready' ? 'btn-outline-primary' : 'btn-default' ">准备</button>
                <button class="btn" @click="onStatus('running')" :class="form.status === 'running' ? 'btn-success' : 'btn-default' ">运行</button>
                <button class="btn" @click="onStatus('stoped')" :class="form.status === 'stoped' ? 'btn-dark' : 'btn-default' ">停止</button>
                <button class="btn" @click="onStatus('fail')" :class="form.status === 'fail' ? 'btn-danger' : 'btn-default' ">异常</button>
            </div>
        </div>

        <div class="col-auto">
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
            onStatus(status){
                this.form.status = status;
                this.onSearch();
            },
            onTag(tag) {
                this.form.tag = tag;
                this.onSearch();
            },
            onNode(node) {
                this.form.node = node;
                this.onSearch();
            },
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
