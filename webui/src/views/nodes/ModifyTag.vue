<template>
    <div class="pull-right">
        <button @click="clickModify" class="btn btn-xs btn-ghost-danger">
            <i class="icon-note"/>
        </button>
        <modal :show.sync="show" title="修改节点标签" @ok="modifyTag" @cancel="editNode=null">
            <div class="form-group">
                <label>标签</label>
                <input class="form-control" type="text" v-model="nodeTag" placeholder="Enter your tag">
            </div>
        </modal>
    </div>
</template>

<script>
    import Modal from "../../plugins/modal";

    export default {
        name: "ModifyTag",
        components: {Modal},
        props: {
            node: Object,
        },
        data: () => ({
            editNode: null,
            nodeTag: "",
        }),
        methods: {
            clickModify() {
                this.editNode = this.node;
                this.nodeTag = this.node.tag;
            },
            modifyTag() {
                var self = this;
                this.editNode = null;
                this.$axios.post("/admin/node/tag", {"key": self.node.key, "tag": self.nodeTag})
                    .then(res => {
                        self.$toast.success("更新成功");
                        self.$emit("change");
                    })
                    .catch(e => {
                        self.$alert("[" + e.error + "]" + e.message);
                    });
                this.$emit("change");
            }
        },
        computed: {
            show() {
                return this.editNode !== null;
            }
        }
    }
</script>

<style scoped>

</style>
