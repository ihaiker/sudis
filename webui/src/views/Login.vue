<template>
    <div class="app flex-row align-items-center bg-dark">
        <vue-particles class="position-absolute w-100 h-100" linesColor="#dedede" color="#dedede"/>

        <div class="container">
            <b-row class="justify-content-center">
                <b-col md="8">
                    <b-card-group>
                        <b-card no-body class="p-4">
                            <b-card-body>
                                <b-form>
                                    <h1 class="text-dark">Login</h1>
                                    <p class="text-muted">Sign In to your account</p>
                                    <b-input-group class="mb-3">
                                        <b-input-group-prepend>
                                            <b-input-group-text><i class="icon-user"></i></b-input-group-text>
                                        </b-input-group-prepend>
                                        <b-form-input type="text" class="form-control" placeholder="Username"
                                                      v-model="name"/>
                                    </b-input-group>
                                    <b-input-group class="mb-4">
                                        <b-input-group-prepend>
                                            <b-input-group-text><i class="icon-lock"></i></b-input-group-text>
                                        </b-input-group-prepend>
                                        <b-form-input type="password" class="form-control" placeholder="Password"
                                                      v-model="passwd"/>
                                    </b-input-group>
                                    <b-row>
                                        <b-col cols="6">
                                            <b-button variant="primary" class="px-4 btn-block" @click="login">Login
                                            </b-button>
                                        </b-col>
                                        <b-col cols="6" class="text-right">
                                            <!--<b-button variant="link" class="px-0">Forgot password?</b-button>-->
                                        </b-col>
                                    </b-row>
                                </b-form>
                            </b-card-body>
                        </b-card>
                        <b-card class="text-white bg-primary py-3 d-md-down-none" style="width:44%">
                            <b-card-body class="text-center">
                                <div>
                                    <img src="@/assets/images/logo2.png" style="width: 3rem;"/>
                                    <h4>Sudis</h4>
                                    <p>Distributed supervisor process control system .</p>
                                    <a href="https://github.com/ihaiker/sudis" target="_blank"
                                       class="btn btn-primary active mt-3">Read More!</a>
                                </div>
                            </b-card-body>
                        </b-card>
                    </b-card-group>
                </b-col>
            </b-row>
        </div>
    </div>
</template>

<script>
    export default {
        name: 'Login',
        data: () => ({
            name: "",
            passwd: "",
        }),
        mounted() {
            let token = localStorage.getItem("token");
            if (token) {
                this.$router.push("/admin/dashboard")
            }
        },
        methods: {
            login() {
                let self = this;
                self.$axios.post("/login", {name: self.name, passwd: self.passwd}).then(res => {
                    let token = res.token;
                    localStorage.setItem("token", token);
                    self.$toast.success('登录成功！');
                    self.$router.push("/admin/dashboard")
                }).catch(e => {
                    self.$alert(e.message);
                });
            }
        }
    }
</script>
