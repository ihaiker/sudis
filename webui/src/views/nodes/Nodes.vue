<template>
    <div>
        <v-title title="节点管理" title-class="icon-puzzle"/>

        <div class="animated fadeIn p-3">

            <table class="table table-hover table-bordered table-fixed  table-striped table-condensed">
                <thead>
                <tr>
                    <th>节点标识</th>
                    <th>节点IP</th>
                    <th>备注</th>
                    <th>地理位置</th>
                    <th style="width: 120px;">管理进程数</th>
                    <th style="width: 100px;">运行状态</th>
                    <th style="width: 180px;">
                        <center>加入日期</center>
                    </th>
                </tr>
                </thead>
                <tbody>
                <tr v-for="(node) in nodes">
                    <td>{{node.key}}</td>
                    <td>{{node.ip}}</td>
                    <td>{{node.tag}}
                        <ModifyTag :node="node" @change="queryNode"/>
                    </td>
                    <td>{{node.address}}</td>
                    <td>{{node.programNum}}</td>
                    <td>
                        <span v-if="node.status === 'online'" class="btn btn-xs btn-success">
                            <i class="fa fa-check"></i> 在线
                        </span>
                        <span v-else class="btn btn-xs btn-danger" @click="forceReload(node)">
                            <i class="fa fa-close"></i> 掉线
                        </span>
                    </td>
                    <td>{{node.time}}</td>
                </tr>
                </tbody>
            </table>
        </div>
    </div>
</template>

<script>

    import ModifyTag from "./ModifyTag";
    import vTitle from "../../plugins/vTitle";

    export default {
        name: "Nodes",
        components: {vTitle, ModifyTag},
        data: () => ({
            nodes: [],
        }),
        mounted() {
            this.queryNode();
        },
        methods: {
            queryNode() {
                this.$axios.get("/admin/node/list", {})
                    .then(res => this.nodes = res)
                    .catch(e => this.$toast.error(e.message))
            },
            forceReload(node) {
                var self = this;
                self.$confirm('删除节点？').then(res => {
                    self.$axios.delete("/admin/node/" + node.key)
                        .then(res => {
                            self.$toast.success("同步成功！");
                            self.queryNode();
                        }).catch(e => {
                        self.$alert("[" + e.error + "]" + e.message);
                    });
                });
            }
        }
    }
</script>
