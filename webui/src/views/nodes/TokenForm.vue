<template>
    <modal :show.sync="show" title="编辑节点TOKEN" @ok="modifyToken" @cancel="cancelModify">
        <div class="form-group">
            <label>标签</label>
            <input class="form-control" type="text" v-model="nodeKey" placeholder="节点唯一码" :disabled="edit">
        </div>
        <div class="form-group">
            <label>
                TOKEN:
                <i class="icon icon-question text-danger"> 此值用于不同节点连接效验盐值</i>
                <button class="btn btn-link" @click="randomToken">
                    <i class="icon icon-reload"></i> 随机一个
                </button>
            </label>
            <input class="form-control" type="text" v-model="nodeToken" placeholder="">
        </div>
    </modal>
</template>

<script>
    import Modal from "../../plugins/modal";

    export default {
        name: "TokenForm",
        components: {Modal},
        props: {
            node: Object,
        },
        data: () => ({
            editNodeToken: null,
            nodeKey: "",
            nodeToken: "",
            edit: false,
        }),
        methods: {
            modifyToken() {
                let self = this;
                self.$axios.post("/admin/node/token", {
                    "key": self.nodeKey, "token": self.nodeToken,
                }).then(res => {
                    self.$toast.success("更新成功");
                    self.$emit("change");
                    self.editNodToken = null;
                }).catch(e => {
                    self.$alert("[" + e.error + "]" + e.message);
                });
            },
            cancelModify() {
                let self = this;
                self.$emit("change");
                self.editNodToken = null;
            },
            randomString(e) {
                e = e || 32;
                let t = "ABCDEFGHJKMNPQRSTWXYZ2345678";
                let str = "";
                for (let i = 0; i < e; i++) {
                    str += t.charAt(Math.floor(Math.random() * t.length));
                }
                return str
            },
            randomToken() {
                this.nodeToken = this.randomString(16);
            }
        },
        watch: {
            node(value) {
                this.edit = false;
                this.nodeKey = "";
                this.randomToken();

                this.editNodeToken = value;
                if (value !== null && JSON.stringify(value) !== "{}") {
                    this.edit = true;
                    this.nodeKey = value.key;
                    this.nodeToken = value.token;
                }
            }
        },
        computed: {
            show() {
                return this.editNodeToken !== null;
            }
        }
    }
</script>
