<template>

    <div class="form-group">
        <label class="col-form-label">
            {{title}} <i class="fa fa-plus text-danger" @click="editArgs.push('')">添加</i>
        </label>
        <div class="controls" v-if="editArgs.length > 0">
            <div v-for="(arg,idx) in editArgs" class="input-group input-group-sm">
                <input class="form-control" v-model="editArgs[idx]" type="text" :placeholder="placeholder">
                <div class="input-group-append">
                    <button class="btn btn-outline-danger" @click="editArgs.splice(idx,1)">
                        <i class="fa fa-remove"/>
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
    export default {
        name: "args",
        props: {
            title: String,
            args: Array,
            placeholder: {
                type: String, default: ""
            }
        },
        mounted() {
            if (this.args) {
                this.editArgs = this.args;
            }
        },
        data: () => ({
            editArgs: [],
        }),
        watch: {
            editArgs(value) {
                this.$emit("change", value);
            }
        }
    }
</script>
