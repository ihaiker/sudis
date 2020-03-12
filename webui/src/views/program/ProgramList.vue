<template>
    <div>
        <v-title title-class="fa-microchip" title="进程管理"/>

        <div class="animated fadeIn p-3">
            <search :nodes="nodes" :tags="tags" @search="queryPrograms">
                <button class="btn btn-sm btn-default ml-2" @click="editProgram={}">
                    <i class="fa fa-plus"/> 新建进程
                </button>
            </search>

            <table class="table table-hover table-bordered table-fixed  table-striped table-condensed">
                <thead>
                <tr>
                    <th>进程名</th>
                    <th>节点</th>
                    <th>标签</th>
                    <th style="width: 285px;">状态 &amp; 操作</th>
                    <th style="width: 55px;">日志</th>
                    <th style="width: 180px;">状态时间</th>
                </tr>
                </thead>
                <tbody>
                <tr v-for="(p) in programs.data">
                    <td>
                        <router-link :to="{path:'/admin/program',query:{node:p.node,name:p.name}}">
                            {{p.name}}
                        </router-link>
                        <template v-if="p.description">
                            <br/><span class="small">{{p.description}}</span>
                        </template>
                    </td>
                    <td>
                        {{nodeShow(p.node)}}
                        <span v-if="!nodeOnline(p.node)" class="text-danger">
                             <span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"/>
                        </span>
                    </td>
                    <td>
                        <Tags :tags="tags" :program="p" @change="queryPrograms"/>
                    </td>
                    <td class="overflow-hidden">
                        <status :program="p" @change="queryPrograms" @edit="editProgram = $event"/>
                    </td>
                    <td>
                        <router-link class="btn btn-sm btn-default" :to="{path:'/admin/program/logs',query:{node:p.node,name:p.name}}">
                            <i class="fa fa-file-text"/>
                        </router-link>
                    </td>
                    <td>{{p.updateTime.substring(0,19).replace("T"," ")}}</td>
                </tr>
                </tbody>
            </table>
            <XPage :items="programs" @change="queryPrograms"/>
        </div>
        <Create :program="editProgram" :nodes="nodes" @change="queryPrograms"/>
    </div>
</template>

<script>
    import Search from "./search";
    import Status from "./status";
    import Create from "./create";
    import Tags from "./tags";
    import vTitle from "../../plugins/vTitle";
    import XPage from "../../plugins/XPage";

    export default {
        name: "ProgramList",
        components: {XPage, vTitle, Tags, Create, Status, Search},
        data: () => ({
            programs: {},
            tags: [], nodes: [],
            editProgram: null, form: {},
        }),
        mounted() {
            this.queryNodes();
            this.queryTags();
        },
        methods: {
            queryNodes() {
                this.$axios.get("/admin/node/list", {})
                    .then(res => this.nodes = res)
                    .catch(e => this.$toast.error(e.error, e.message))
            },
            queryTags() {
                let self = this;
                this.request("加载标签...", this.$axios.get("/admin/tag/list").then(res => {
                    self.tags = res;
                }));
            },

            queryPrograms(form) {
                this.form = this.twoJsonMerge(this.form, form || {});
                let self = this;
                self.programs = {};
                this.$axios.get("/admin/program/list?" + this.$form.transformRequest[0](this.form))
                    .then(res => self.programs = res)
                    .catch(e => {
                        self.$toast.error(e.message);
                    })
            },

            nodeShow(nodeKey) {
                for (let idx in this.nodes) {
                    let node = this.nodes[idx];
                    if (node.key === nodeKey) {
                        return node.tag === "" ? node.key : node.tag;
                    }
                }
                return nodeKey;
            },
            nodeOnline(nodeKey) {
                for (let idx in this.nodes) {
                    let node = this.nodes[idx];
                    if (node.key === nodeKey) {
                        return node.status === 'online';
                    }
                }
                return false;
            }
        }
    }
</script>
<style>
    .spinner-border-xs {
        width: 0.65rem;
        height: 0.65rem;
        border-width: 0.2em;
    }
    table.table-hover>tbody>tr:hover>td, table.table-hover>tbody>tr:hover>th {
        background-color: rgba(37, 236, 49, 0.29);
    }
</style>
