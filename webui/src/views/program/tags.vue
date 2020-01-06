<template>
    <div>
        <button v-if="program.tags && program.tags.length > 0" v-for="tag in program.tags" @click="onEdit"
                :class="tagShow(tag)" class="btn btn-xs mr-1">{{tag}}
        </button>

        <button v-if="!program.tags || program.tags.length === 0" class="btn btn-xs btn-default" @click="onEdit">
            <i class="fa fa-plus"/>
        </button>

        <modal :show.sync="show" title="编辑标签" @cancel="onClickOK" @ok="onClickOK">
            <div class="b-b-1 mb-3">系统定义标签</div>
            <div class="row">
                <div v-for="tag in tags" class="col-auto mb-3">
                    <div class="input-group">
                        <div class="input-group-prepend">
                            <button :class="tag.class" class="btn btn-sm">{{tag.name}}</button>
                        </div>
                        <c-switch color="primary" :checked="hasTag(tag.name)" @change="submitEdit(tag.name,$event)"/>
                    </div>
                </div>
            </div>

            <div class="b-b-1 mb-3">用户标签</div>
            <div class="row">
                <div v-for="tag in program.tags" v-if="tagShow(tag) === 'btn-outline-primary'" class="col-auto mb-3">
                    <div class="input-group">
                        <div class="input-group-prepend">
                            <button class="btn btn-sm btn-outline-primary">{{tag}}</button>
                        </div>
                        <c-switch class="switch-lg" color="primary" :checked="true" @change="submitEdit(tag,$event)"/>
                    </div>
                </div>
            </div>
            <div class="input-group input-group-sm">
                <input class="form-control" type="text" v-model="tagName" placeholder="标签内容"/>
                <div class="input-group-append">
                    <button class="btn btn-primary" type="button" @click="submitEdit(tagName,true)">
                        <i class="fa fa-plus"/> 添加
                    </button>
                </div>
            </div>
        </modal>
    </div>
</template>

<script>
    import Modal from "../../plugins/modal";
    import {Switch as cSwitch} from '@coreui/vue'

    export default {
        name: "tags",
        components: {Modal, cSwitch},
        props: {
            program: Object,
            tags: Array
        },
        data: () => ({
            show: false,
            tagName: "",
        }),
        methods: {
            hasTag(tagName) {
                for (let idx in this.program.tags) {
                    if (this.program.tags[idx] === tagName) {
                        return true;
                    }
                }
                return false;
            },
            tagShow(tagName) {
                for (let idx in this.tags) {
                    let tag = this.tags[idx];
                    if (tag.name === tagName) {
                        return tag.class;
                    }
                }
                return "btn-outline-primary";
            },
            onEdit() {
                this.show = true;
            },
            onClickOK(){
                this.show = false;
                this.$emit("change");
            },
            submitEdit(tag, checked) {
                this.tagName = "";
                if (this.hasTag(tag) && checked) {
                    return;
                } else if (!this.hasTag(tag) && !checked) {
                    return;
                }
                let self = this;
                let params = {
                    name: self.program.name, node: self.program.node,
                    tag: tag, add: (checked ? 1 : 0)
                };
                self.$axios.post("/admin/program/tag", params).then(res => {
                    if (checked) {
                        self.program.tags.push(tag);
                    } else {
                        let i = self.program.tags.indexOf(tag);
                        self.program.tags.splice(i, 1);
                    }
                }).catch(e => {
                    self.$toast.error(e.error, e.message);
                });
            }
        }
    }
</script>
