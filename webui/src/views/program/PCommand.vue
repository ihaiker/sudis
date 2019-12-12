<template>
    <div>
        <div class="form-group" :class="{'mb-0':daemon}">
            <label class="col-form-label">{{title}}命令</label>
            <div class="controls">
                <input class="form-control form-control-sm" v-model="commandName" :placeholder="title + '程序'">
            </div>
            <args class="mb-0" :args="commandArgs" :title="title + '参数'" @change="commandArgs = $event"/>
        </div>

        <div v-if="daemon" class="form-group">
            <label class="col-form-label">状态检查</label>
            <div class="input-prepend input-group input-group-sm">
                <select v-model="healthType" class="form-control" style="width: 60px !important; flex: none;">
                    <option value="http">http</option>
                    <option value="https">https</option>
                    <option value="tcp">tcp</option>
                </select>
                <input v-model="healthUrl" class="form-control" placeholder="地址">
            </div>
            <div class="input-group input-group-sm">
                <div class="input-group-prepend">
                    <span class="input-group-text">
                        检查周期
                    </span>
                </div>
                <input v-model="healthTTL" class="form-control" type="number" placeholder="TTL(秒)" style="width: 80px !important; flex: none;">
                <div class="input-group-prepend">
                    <span class="input-group-text">
                        安全码
                    </span>
                </div>
                <input v-model="healthSecurityKey" class="form-control" placeholder="security token">
            </div>
        </div>
    </div>
</template>

<script>
    import Args from "./args";

    export default {
        name: "vCommand",
        components: {Args},
        props: {
            title: String,
            daemon: Boolean,
            command: Object,
        },
        mounted() {
            this.commandName = this.command.command;
            if (this.command.args) {
                this.commandArgs.push(...this.command.args);
            }
            if (this.command.health) {
                this.healthType = this.command.health.type;
                this.healthUrl = this.command.health.url;
                this.healthTTL = this.command.health.ttl;
                this.healthSecurityKey = this.command.health.securityKey;
            }
        },
        data: () => ({
            commandName: "",
            commandArgs: [],
            healthType: "http",
            healthUrl: "",
            healthTTL: 5,
            healthSecurityKey: ""
        }),
        methods: {
            notify() {
                let out = {
                    command: this.commandName, args: this.commandArgs,
                };
                if (this.daemon) {
                    out.health = {
                        type: this.healthType, url: this.healthUrl,
                        ttl: this.healthTTL, securityKey: this.healthSecurityKey
                    }
                }
                this.$emit("change", out);
            }
        },
        watch: {
            //@formatter:off
            commandName(value) {
                let args = value.split(" ");
                if(args.length > 1){
                    for(let i = 1; i < args.length; i++) {
                        if(args[i] !== ""){
                            this.commandArgs.push(args[i]);
                        }
                    }
                    this.commandName = args[0];
                }
                this.notify();
            },
            commandArgs(args) {this.notify();},
            healthType() { this.notify(); },
            healthUrl(){ this.notify(); },
            healthTTL(){ this.notify(); },
            healthSecurityKey(){this.notify();}
            //@formatter:on
        }
    }
</script>
