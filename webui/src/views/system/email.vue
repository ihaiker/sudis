<template>
    <div>
        <div class="alert alert-info">邮件服务器配置</div>
        <div class="row ">
            <div class="col-4 form-group">
                <div class="input-group">
                    <div class="input-group-prepend">
                        <span class="input-group-text">
                            <i class="fa fa-user"/> SMTP服务地址：
                        </span>
                    </div>
                    <input class="form-control" v-model="address" type="text" placeholder="SMTP地址"/>
                </div>
            </div>
            <div class="col-4 form-group">
                <div class="input-group">
                    <div class="input-group-prepend">
                        <span class="input-group-text">
                            <i class="fa fa-envelope"/>
                        </span>
                    </div>
                    <input class="form-control" v-model="port" type="number" placeholder="SMTP端口号"/>
                </div>
            </div>
            <div class="col-4">
                <button class="btn btn-primary w-25">设置</button>
            </div>
            <div class="col"></div>
        </div>

        <div class="row">
            <div class="col form-group">
                <div class="input-group">
                    <div class="input-group-prepend">
                        <span class="input-group-text">
                            <i class="fa fa-user"/>&nbsp;用户名
                        </span>
                    </div>
                    <input class="form-control" v-model="name" type="text" placeholder="Username">
                </div>
            </div>
            <div class="col form-group">
                <div class="input-group">
                    <div class="input-group-prepend">
                        <span class="input-group-text">
                            <i class="fa fa-asterisk"/>&nbsp;密码
                        </span>
                    </div>
                    <input class="form-control" v-model="passwd" type="password" placeholder="Password">
                </div>
            </div>
            <div class="col"></div>
        </div>
        <div class="alert alert-info mt-2">通知模板</div>
        <pre> <span v-for="(name,attr) in attrs">{{attr}}：{{name}}  </span></pre>
        <textarea class="form-control mt-1" v-model="content"></textarea>
    </div>
</template>

<script>
    export default {
        name: "email",
        data: () => ({
            attrs: {
                Id: "程序ID", Name: "程序名称", WorkDir: "工作目录", User: "运行用户",
                Envs: "环境变量", State: "状态",
                "Start.Command": "启动命令", "Start.Args": "启动命令参数",
                "Stop.Command": "停止命令", "Stop.Args": "停止命令参数",
            },
            address: "", port: 465,
            name: "", passwd: "",
            content: `节点：{{.Node}}，程序：{{.Name}}，状态更改：{{.State}}`,
        })
    }
</script>
