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
                    <th style="width: 280px;">状态 &amp; 操作</th>
                    <th style="width: 180px;">状态时间</th>
                </tr>
                </thead>
                <tbody>
                <tr v-for="(p) in programs">
                    <td>{{p.name}}</td>
                    <td>{{nodeShow(p.node)}}</td>
                    <td>
                        <Tags :tags="tags" :program="p" @change="queryPrograms"/>
                    </td>
                    <td>
                        <status :program="p" @change="queryPrograms" @edit="editProgram = $event"/>
                    </td>
                    <td>{{p.time}}</td>
                </tr>
                </tbody>
            </table>
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

    export default {
        name: "Programs",
        components: {vTitle, Tags, Create, Status, Search},
        data: () => ({
            programs: [],
            tags: [],
            nodes: [],
            editProgram: null
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
                let requestData = form || {};
                let self = this;
                self.programs = [];
                this.$axios.get("/admin/program/list?" + this.$form.transformRequest[0](requestData))
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
            }
        }
    }
</script>
