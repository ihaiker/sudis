<template>
    <modal :show.sync="show" title="标签管理" @ok="ok" @cancel="cancel">
        <div class="form-horizontal">
            <div class="form-group row">
                <div class="col-md-12">
                    <div class="input-group">
                        <span class="input-group-prepend">
                            <button class="btn btn-primary" type="button">
                                <i class="fa fa-search"></i> 标签名称</button>
                            </span>
                        <input class="form-control" type="text" v-model="formTag.name" placeholder="标签名称：只能字母、数字和下划线">
                    </div>
                </div>
            </div>
        </div>

        <div class="form-group">
            <label for="name">标签样式：</label>
            <div class="controls">
                <button v-for="c in buttons" :class="c" class="btn btn-sm ml-3 mt-3" @click="formTag.class=c">{{formTag.name}}</button>
            </div>
        </div>

    </modal>
</template>

<script>
    import Modal from "../../plugins/modal";

    export default {
        name: "TagForm",
        props: {
            tag: Object,
        },
        data: () => ({
            show: false,
            buttons: [
                "btn-pinterest", "btn-linkedin", "btn-yahoo",
                "btn-primary", "btn-secondary", "btn-success", "btn-danger", "btn-warning", "btn-info", "btn-dark",
                "btn-outline-primary", "btn-outline-secondary", "btn-outline-success", "btn-outline-danger",
                "btn-outline-warning", "btn-outline-info", "btn-outline-dark",

                "btn-pill btn-pinterest", "btn-pill btn-linkedin", "btn-pill btn-yahoo",
                "btn-pill btn-primary", "btn-pill btn-success", "btn-pill btn-danger",
                "btn-pill btn-warning", "btn-pill btn-info", "btn-pill btn-dark",
                "btn-pill btn-outline-primary", "btn-pill btn-outline-secondary", "btn-pill btn-outline-success", "btn-pill btn-outline-danger",
                "btn-pill btn-outline-warning", "btn-pill btn-outline-info", , "btn-pill btn-outline-dark",
            ],
            formTag: {name: "", "class": ""}
        }),
        components: {Modal},
        methods: {
            ok() {
                var self = this;
                this.request("编辑标签中...",
                    this.$axios.post("/admin/tag/addOrModify", this.formTag).then(res => {
                        self.show = false;
                        self.$emit("ok");
                    }).catch(e => {
                        self.$alert(e.message);
                    })
                );
            },
            cancel() {
                this.show = false;
                this.$emit("cancel");
            }
        },
        watch: {
            tag(value) {
                if (value) {
                    this.show = true;
                    this.formTag = value;
                } else {
                    this.show = false;
                    this.formTag = {name: "", "class": ""};
                }
            }
        },
    }
</script>

<style scoped>

</style>
