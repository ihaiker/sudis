<template>
    <div>
        <button class="btn btn-sm ml-1" :class="statusColor(wProgram.status)" style="width: 70px;">
            {{wProgram.status}}
        </button>

        <template v-if="isRunning">
            <button v-if="isStopting" class="btn btn-sm btn-default ml-1" type="button" disabled>
                <span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"/>&nbsp;停止
            </button>
            <button v-else @click="stopCommand" class="btn btn-sm btn-default ml-1" :disabled="disable">
                <i class="fa fa-stop">&nbsp;停止</i>
            </button>

            <button class="btn btn-sm btn-default ml-1" @click="restartCommand" :disabled="disable">
                <i class="fa fa-refresh">&nbsp;重启</i>
            </button>
            <router-link :to="{path:'/admin/program',query:{node:wProgram.node,name:wProgram.name}}"
                         class="btn btn-sm btn-default ml-1" :disabled="disable">
                <i class="fa fa-info-circle">&nbsp;详情</i>
            </router-link>
        </template>
        <template v-else>
            <button v-if="isStarting" class="btn btn-sm btn-default ml-1" type="button" disabled>
                <span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"/>&nbsp; 启动
            </button>
            <button v-else @click="startCommand" class="btn btn-sm btn-default ml-1" :disabled="disable">
                <i class="fa fa-play">&nbsp;启动</i>
            </button>

            <button class="btn btn-sm btn-default ml-1" :disabled="disable" @click="onEditClick">
                <i class="fa fa-edit">&nbsp;编辑</i>
            </button>

            <button @click="deleteCommand" class="btn btn-sm btn-default ml-1" :disabled="disable">
                <i class="fa fa-trash">&nbsp;删除</i>
            </button>
        </template>
    </div>
</template>

<script>
    let Starting = "starting";
    let Running = "running";
    let Fail = "fail";
    let Ready = "ready";
    let Stoped = "stoped";
    let Stoping = "stoping";

    export default {
        name: "status",
        props: {
            program: Object,
        },
        data: () => ({
            wProgram: null,
            disable: false,
        }),
        created() {
            this.wProgram = this.program;
            this.disable = (this.program.status === Starting || this.program.status === Stoping);
        },
        methods: {
            statusColor(status) {
                if (status === Running || status === Starting) {
                    return "btn-success";
                } else if (status === Fail) {
                    return "btn-pinterest"
                }
                return "btn-default";
            },
            onChange() {
                let self = this;
                setTimeout(() => {
                    self.$emit("change");
                }, 1000)
            },
            beforeCommand(status) {
                this.disable = true;
                this.wProgram.status = status;
            },

            afterCommand(status, err) {
                this.disable = false;
                this.wProgram.status = status;
                if (err != null) {
                    this.$toast.error(err.message);
                }
            },

            startCommand() {
                let self = this;
                self.beforeCommand(Starting);
                self.execCommand("start", res => {
                    self.afterCommand(Running)
                }, e => {
                    self.afterCommand(Fail, e)
                })
            },

            restartCommand() {
                let self = this;
                self.beforeCommand(Stoping);
                self.execCommand("restart", res => {
                    self.afterCommand(Running)
                }, e => {
                    self.afterCommand(Fail, e)
                })
            },

            stopCommand() {
                let self = this;
                self.beforeCommand(Stoping);
                self.execCommand("stop", res => {
                    self.afterCommand(Stoped);
                }, e => {
                    self.afterCommand(Fail, e);
                })
            },

            deleteCommand() {
                let self = this;
                self.$confirm('确定删除吗？').then(res => {
                    self.beforeCommand(Stoped);
                    self.execCommand("delete", () => {
                        self.onChange()
                    }, (e) => {
                        self.afterCommand(Stoped, e);
                    })
                });
            },

            execCommand(command, ok, err, fin) {
                let params = {
                    name: this.wProgram.name,
                    node: this.wProgram.node,
                    command
                };
                let self = this;
                self.$axios.put("/admin/program/command", params)
                    .then(ok).catch(err).finally(fin);
            },
            onEditClick() {
                let self = this;
                let params = this.$form.transformRequest[0]({
                    name: this.wProgram.name,
                    node: this.wProgram.node,
                });
                self.$axios.get("/admin/program/detail?" + params).then(res => {
                    self.$emit("edit", res);
                }).catch(e => {
                    self.$toast.error(e.error, e.message);
                });
            }
        },
        computed: {
            isRunning() {
                return (this.wProgram.status === Running || this.wProgram.status === Stoping)
            },
            isStarting() {
                return this.wProgram.status === Starting
            },
            isStopting() {
                return this.wProgram.status === Stoping
            }
        }
    }
</script>
