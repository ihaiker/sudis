<template>
    <div>
        <v-title title="节点管理" title-class="icon-puzzle"/>

        <div class="animated fadeIn p-3">

            <div class="form-group">
                <div class="input-group">
                    <div class="input-group-prepend">
                        <span class="input-group-text">节点名称</span>
                    </div>
                    <input class="form-control" type="text" name="program" placeholder="节点名称"/>
                    <div class="input-group-prepend">
                        <button class="btn btn-danger" @click="editNodToken={}">
                            <i class="fa fa-plus"/> 添加Token
                        </button>
                    </div>
                </div>
            </div>

            <table class="table table-hover table-bordered table-fixed  table-striped table-condensed mt-2">
                <thead>
                <tr>
                    <th>节点标识(节点IP)</th><th>Token</th><th>备注</th><th>地理位置</th>
                    <th style="width: 70px;">进程数</th>
                    <th style="width: 90px;">运行状态</th>
                    <th style="width: 180px;">
                        <center>加入日期</center>
                    </th>
                </tr>
                </thead>
                <tbody>
                <tr v-for="(node) in nodes">
                    <td>{{node.key}} ({{node.ip}})</td>
                    <td>
                        {{node.token}}
                        <button @click="editNodToken=node" class="btn btn-xs btn-ghost-danger pull-right" >
                            <i class="icon-note"/>
                        </button>
                    </td>
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

            <TokenForm :node="editNodToken" @change="queryNode"/>
        </div>
    </div>
</template>

<script>

    import ModifyTag from "./ModifyTag";
    import vTitle from "../../plugins/vTitle";
    import TokenForm from "./TokenForm";

    export default {
        name: "Nodes",
        components: {TokenForm, vTitle, ModifyTag},
        data: () => ({
            nodes: [],
            editNodToken: null,
        }),
        mounted() {
            this.queryNode();
        },
        methods: {
            queryNode() {
                this.editNodToken = null;
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
