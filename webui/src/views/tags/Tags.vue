<template>
    <div>
        <v-title title="标签管理" title-class="fa-tags"/>

        <div class="animated fadeIn p-2">
            <button @click="editTag={}" class="btn btn-pinterest">
                <i class="icon-plus"/> &nbsp;添加标签
            </button>
        </div>

        <div class="animated fadeIn p-1 pl-3 pr-3">
            <div class="row text-center">
                <div v-for="(tag) in tags" class="col-1">
                    <div class="row mt-3 ml-1">

                        <router-link :to="'/admin/programs/' + tag.name " class="col-12 btn btn-block" :class="tag.class">{{tag.name}}</router-link>

                        <button @click="editTag=tag" class="col-6 btn btn-outline-primary btn-block mt-0">
                            <i class="fa fa-edit"/>
                        </button>

                        <button class="col-6 btn btn-outline-danger btn-block mt-0">
                            <delete @ok="removeTag(tag.name)">
                                <i class="fa fa-trash"/>
                            </delete>
                        </button>
                    </div>
                </div>
            </div>
        </div>

        <tag-form :tag="editTag" @ok="queryTags" @cancel="queryTags"/>
    </div>
</template>

<script>
    import TagForm from "./TagForm";
    import Delete from "../../plugins/delete";
    import vTitle from "../../plugins/vTitle";

    export default {
        name: "Tags",
        components: {vTitle, Delete, TagForm},
        data: () => ({
            tags: [],
            editTag: null,
        }),
        mounted() {
            this.queryTags();
        },
        methods: {
            queryTags() {
                var self = this;
                self.editTag = null;
                this.request("加载标签...", this.$axios.get("/admin/tag/list").then(res => {
                    self.tags = res;
                }).catch(e => self.$toast.error(e.message)));
            },
            removeTag(name) {
                var self = this;
                this.request("删除标签...",
                    this.$axios.delete("/admin/tag/" + name)
                        .then(res => {
                            self.queryTags();
                            self.$toast.success("删除标签" + name + "成功!")
                        }).catch(e => self.$toast.error(e.error, e.message))
                );
            }
        }
    }
</script>
