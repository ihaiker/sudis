<template>
    <modal v-if="show" :show.sync="show" :title="title" @ok="ok" @cancel="cancel">
        <div class="form-group">
            <div class="input-group input-group-sm">
                <div class="input-group-prepend">
                    <div class="input-group-text">程序名称</div>
                </div>
                <input class="form-control" v-model="form.name" placeholder="只能为字母、数字和下划线组合！" :disabled="edit">
            </div>
            <div class="input-group input-group-sm mt-2">
                    <span class="input-group-prepend">
                        <span class="input-group-text">所属节点</span>
                    </span>
                <select class="form-control" v-model="form.node" :disabled="edit">
                    <option value="">选择节点</option>
                    <option v-for="node in nodes" :value="node.key">{{node.tag === "" ? node.key : node.tag}}</option>
                </select>
            </div>
            <div class="input-group input-group-sm mt-2">
                <div class="input-group-prepend">
                    <div class="input-group-text">描述信息</div>
                </div>
                <input class="form-control" v-model="form.description" placeholder="描述信息">
            </div>

            <div class="input-group input-group-sm mt-2">
                <div class="input-group-prepend">
                    <span class="input-group-text">进程模式</span>
                </div>
                <select class="form-control" v-model="form.daemon">
                    <option value="0">非daemon进程</option>
                    <option value="1">daemon进程</option>
                </select>
            </div>
        </div>

        <PCommand title="启动" :daemon="form.daemon === '1'" :command="form.start" @change="form.start = $event"/>
        <PCommand v-if="form.daemon === '1'" title="停止" :daemon="false" :command="form.stop"
                  @change="form.stop = $event"/>

        <div class="form-group">
            <div class="input-group input-group-sm">
                <div class="input-group-prepend">
                    <span class="input-group-text">日志文件</span>
                </div>
                <input class="form-control" v-model="form.logger" placeholder="日志文件">
            </div>
        </div>

        <div class="form-group">
            <div class="input-group input-group-sm">
                <div class="input-group-prepend">
                    <span class="input-group-text">工作目录</span>
                </div>
                <input class="form-control" v-model="form.workDir" placeholder="默认为用户主目录">
            </div>

            <div class="input-group input-group-sm mt-2">
                <div class="input-group-prepend">
                    <span class="input-group-text">启动用户</span>
                </div>
                <input class="form-control" v-model="form.user" placeholder="默认为主程序启动用户">
            </div>
        </div>

        <args :args="form.envs" title="环境变量" placeholder="环境变量：填写NAME=VALUE" @change="form.envs = $event"/>

        <div class="form-group row">
            <div class="col-auto col-form-label mr-0 pr-0">
                自动启动：
            </div>
            <div class="col-auto col-form-label ml-0 pl-0">
                <cSwitch class="switch-sm" color="primary" v-model="form.autoStart"/>
            </div>

            <div class="col-auto col-form-label" v-if="form.daemon==='1'">
                忽略已启动
            </div>
            <div class="col-auto col-form-label ml-0 pl-0" v-if="form.daemon==='1'">
                <cSwitch class="switch-sm" color="primary" v-model="form.ignoreStarted"/>
            </div>
        </div>

        <div class="form-group row">
            <div class="col-auto  col-form-label mr-0 pr-0">
                启动预计时间：
            </div>
            <div class="col-2 ml-0 pl-0">
                <input class="form-control form-control-sm" type="number" v-model="form.startDuration"
                       placeholder="(秒)"/>
            </div>
            <div class="col-auto col-form-label mr-0 pr-0">
                失败重试次数：
            </div>
            <div class="col-2 ml-0 pl-0">
                <input class="form-control form-control-sm" type="number" v-model="form.startRetries"
                       placeholder="(次)"/>
            </div>
        </div>

    </modal>
</template>

<script>
    import Modal from "../../plugins/modal";
    import PCommand from "./PCommand";
    import {Switch as cSwitch} from '@coreui/vue'
    import Args from "./args";

    export default {
        name: "create", components: {Args, cSwitch, PCommand, Modal},
        props: {
            program: {
                type: Object,
                default: null,
            },
            nodes: {
                type: Array,
                default: [],
            }
        },
        data: () => ({
            form: null, title: "", edit: false,
        }),
        methods: {
            onChange() {
                let self = this;
                setTimeout(() => {
                    self.$emit("change");
                }, 1000)
            },
            ok() {
                let self = this;
                self.form.startDuration = parseInt(self.form.startDuration);
                self.form.startRetries = parseInt(self.form.startRetries);
                self.$axios.post("/admin/program/addOrModify", self.form).then(res => {
                    self.$toast.success("添加成功！");
                    self.form = null;
                    self.onChange();
                }).catch(e => {
                    self.$alert(e.message);
                });
            },
            cancel() {
                this.form = null;
            },
            reset() {
                this.form = {
                    name: "", node: "", daemon: "0",
                    autoStart: false, ignoreStarted: false,
                    startDuration: 3, startRetries: 3, logger: "",
                    envs: [],
                    start: {
                        command: "", args: [],
                        health: {
                            type: "", url: "", ttl: 5, securityKey: ""
                        }
                    },
                    stop: {
                        command: "", args: [],
                    },
                    workDir: "", user: "",
                };
            }
        },
        computed: {
            show() {
                return this.form !== null;
            }
        },
        watch: {
            program(value) {
                if (JSON.stringify(value) === "{}") {
                    this.edit = false;
                    this.title = "添加程序";
                    this.reset();
                } else {
                    this.edit = true;
                    this.title = "编辑：" + value.name;
                    this.form = value;
                }
            }
        }
    }
</script>
