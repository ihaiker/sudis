<template>
    <div>
        <v-title title="用户管理" title-class="fa fa-user">
            <button class="btn btn-xs btn-pinterest" @click="addUser={'user':'','passwd':''}">
                <i class="icon-plus"/> &nbsp;添加用户
            </button>
        </v-title>
        <div class="animated fadeIn p-5">
            <div class="row">
                <table class="col-6 table table-hover table-bordered table-fixed table-striped text-center">
                    <thead>
                    <tr>
                        <th>账号</th>
                        <th>时间</th>
                        <th style="width: 200px;">操作</th>
                    </tr>
                    </thead>
                    <tbody>
                    <tr v-for="(user) in users">
                        <td>{{user.name}}</td>
                        <td>{{user.time}}</td>
                        <td>
                            <button class="btn btn-sm btn-primary"
                                    @click="modifyPasswd = {name: user.name, passwd: '',}">
                                <i class="fa fa-refresh"/>&nbsp;重置密码
                            </button>
                            <button class="btn btn-sm btn-danger ml-1">
                                <delete message="您确定要删除用户" @ok="onDeleteUser(user.name)">
                                    <i class="fa fa-remove"/>&nbsp;删除
                                </delete>
                            </button>
                        </td>
                    </tr>
                    </tbody>
                </table>
                <div class="col-6">
                    <div class="card" v-if="addUser !== null">
                        <div class="card-header">添加用户</div>
                        <div class="card-body">
                            <div class="form-group">
                                <div class="input-group">
                                    <div class="input-group-prepend"><span class="input-group-text"><i class="fa fa-user"></i></span></div>
                                    <input class="form-control" type="text" v-model="addUser.name" placeholder="Username" autocomplete="name">
                                </div>
                            </div>
                            <div class="form-group">
                                <div class="input-group">
                                    <div class="input-group-prepend"><span class="input-group-text"><i class="fa fa-asterisk"></i></span></div>
                                    <input class="form-control" type="password" v-model="addUser.passwd" placeholder="Password" autocomplete="new-password">
                                </div>
                            </div>
                        </div>
                        <div class="card-footer">
                            <div class="row">
                                <div class="col-3 offset-6">
                                    <button class="btn btn-sm btn-block btn-light" @click="addUser = null">取消</button>
                                </div>
                                <div class="col-3">
                                    <button class="btn btn-sm btn-block btn-success" @click="onAddUser">确定</button>
                                </div>
                            </div>
                        </div>
                    </div>

                    <div class="card" v-if="modifyPasswd !== null">
                        <div class="card-header">修改密码</div>
                        <div class="card-body">
                            <h4>修改{{modifyPasswd.name}}密码</h4>
                            <div class="form-group">
                                <div class="input-group">
                                    <div class="input-group-prepend"><span class="input-group-text"><i class="fa fa-asterisk"></i></span></div>
                                    <input class="form-control" type="text" v-model="modifyPasswd.passwd" placeholder="新密码" autocomplete="password">
                                </div>
                            </div>
                        </div>
                        <div class="card-footer">
                            <div class="row">
                                <div class="col-3 offset-6">
                                    <button class="btn btn-sm btn-block btn-light" @click="modifyPasswd = null">取消</button>
                                </div>
                                <div class="col-3">
                                    <button class="btn btn-sm btn-block btn-pinterest" @click="onModifyPasswd">重置</button>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
    import VTitle from "../../plugins/vTitle";
    import Delete from "../../plugins/delete";

    export default {
        name: "Users",
        components: {VTitle, Delete},
        data: () => ({
            users: [],
            addUser: null,
            modifyPasswd: null,
        }),
        mounted() {
            this.queryUser();
        },
        methods: {
            queryUser() {
                let self = this;
                self.$axios.get("/admin/user/list").then(res => {
                    self.users = res;
                });
            },
            onAddUser() {
                let self = this;
                self.$axios.post("/admin/user/add", self.addUser).then(res => {
                    self.$toast.success('添加' + self.addUser.name + '用户成功！');
                    self.addUser = null;
                    self.queryUser();
                }).catch(e => {
                    self.$alert(e.message);
                })
            },
            onModifyPasswd() {
                let self = this;
                self.$axios.post("/admin/user/passwd", self.modifyPasswd).then(res => {
                    self.$toast.success('修改' + self.modifyPasswd.name + '密码成功！');
                    self.modifyPasswd = null;
                    self.queryUser();
                }).catch(e => {
                    self.$alert(e.message);
                })
            },
            onDeleteUser(name) {
                let self = this;
                self.$axios.delete("/admin/user/" + name).then(res => {
                    self.$toast.success('删除' + name + '用户成功！');
                    self.queryUser();
                }).catch(e => {
                    self.$alert(e.message);
                })
            }
        }
    }
</script>
