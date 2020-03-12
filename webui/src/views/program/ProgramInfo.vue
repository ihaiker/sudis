<template>
    <div>
        <ol class="breadcrumb breadcrumb-fixed">
            <li class="breadcrumb-item">
                <router-link to="/admin/programs">
                    <i class="fa fa-microchip"/> 进程管理
                </router-link>
            </li>
            <li class="breadcrumb-item">
                {{$route.query.name}}
            </li>
        </ol>
        <div v-if="program" class="p-3">
            <table class="table table-hover table-bordered">
                <tbody>
                <tr>
                    <td width="100px">程序名称：</td>
                    <td><span class="text-primary">{{program.name}}</span> {{program.description}}</td>
                    <td width="100px">所属节点：</td>
                    <td>{{program.node}}</td>
                </tr>
                <tr>
                    <td>启动命令：</td>
                    <td colspan="3">
                        <span class="text-danger">{{program.start.command}}</span>&nbsp;
                        <span v-if="program.start.args">{{program.start.args.join("&nbsp;")}}</span>
                    </td>
                </tr>
                <tr v-if="program.stop">
                    <td>停止命令</td>
                    <td colspan="3">
                        <span class="text-danger">{{program.stop.command}}</span>&nbsp;
                        <span v-if="program.stop.args">{{program.stop.args.join("&nbsp;")}}</span>
                    </td>
                </tr>
                <tr>
                    <td>工作目录</td>
                    <td>{{program.workDir}}</td>
                    <td>日志输出</td>
                    <td>{{program.logger}}</td>
                </tr>
                <tr>
                    <td>启动用户</td>
                    <td>{{program.user}}</td>
                    <td>进程ID</td>
                    <td>{{program.pid}}</td>
                </tr>
                <tr>
                    <td>添加时间</td>
                    <td colspan="3">{{program.addTime}}</td>
                </tr>
                <tr>
                    <td>内存使用</td>
                    <td>{{ramShow(program.rss)}}</td>
                    <td>CPU使用</td>
                    <td>{{program.cpu}}</td>
                </tr>
                </tbody>
            </table>

            <div v-if="program.pid != 0" class="card card-accent-primary">
                <div class="card-header">
                    <i class="icon-speedometer"/>内存使用情况
                </div>
                <RAM :line-data="ramData"/>
            </div>

            <div v-if="program.pid !== 0" class="card card-accent-dark">
                <div class="card-header">
                    <i class="icon-speedometer"/>CPU使用情况
                </div>
                <CPU :line-data="cpuData"/>
            </div>

        </div>
    </div>
</template>

<script>
    import VTitle from "../../plugins/vTitle";
    import RAM from "./RAM";
    import CPU from "./CPU";

    export default {
        name: "ProgramInfo",
        components: {CPU, RAM, VTitle},
        data: () => ({
            _timer: null,
            program: null,
            cpuData: [],
            ramData: [],
        }),
        mounted() {
            this.queryDetail();
            for (let i = 0; i < 200; i++) {
                this.cpuData.push({name: "", value: [i, 0]});
                this.ramData.push({name: "", value: [i, 0]});
            }
        },
        methods: {
            queryDetail() {
                let self = this;
                let params = this.$form.transformRequest[0]({
                    name: this.$route.query.name,
                    node: this.$route.query.node,
                });
                self.$axios.get("/admin/program/detail?" + params).then(res => {
                    self.program = res;
                    if (self.cpuData.length >= 200) {
                        self.cpuData.shift()
                    }
                    if (self.ramData.length >= 200) {
                        self.ramData.shift();
                    }
                    self.show = (self.program.pid !== 0);
                    let now = self.now();
                    self.cpuData.push({name: now.name, value: [now.show, self.program.cpu]});
                    self.ramData.push({name: now.name, value: [now.show, self.program.rss]});
                }).catch(e => {
                    if (self.show) {
                        self.$toast.error(e.error, e.message);
                    }
                }).finally(() => {
                    self._timer = setTimeout(self.queryDetail, 1000)
                });
            }
        },
        beforeDestroy() {
            this._timer && clearTimeout(this._timer);
        }
    }
</script>
