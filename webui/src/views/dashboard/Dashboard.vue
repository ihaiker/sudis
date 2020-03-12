<template>
    <div class="animated fadeIn p-lg-5 p-sm-2">
        <b-row>
            <b-col v-if="volumns.version.tag_name" md="3" sm="6">
                <div class="brand-card">
                    <div class="brand-card-header bg-google-plus">
                        <i class="fa fa-download"/>
                        <div class="chart-wrapper">
                            <social-box-chart :data="[35, 23, 56, 22, 97, 23, 64]"/>
                        </div>
                    </div>
                    <div class="brand-card-body">
                        <a :href="volumns.version.html_url" target="_blank">
                            <div class="text-value text-capitalize">{{volumns.version.tag_name}}</div>
                            <div class="text-uppercase text-muted">最新版本</div>
                        </a>
                        <a class="align-self-center" :href="volumns.version.html_url" target="_blank">
                            <div class="text-uppercase text-muted">{{volumns.version.name}}</div>
                        </a>
                    </div>
                </div>
            </b-col>

            <b-col md="3" sm="6">
                <div class="brand-card">
                    <div class="brand-card-header bg-facebook">
                        <i class="icon-puzzle"> Nodes</i>
                        <div class="chart-wrapper">
                            <social-box-chart :data="[65, 59, 84, 84, 51, 55, 40]"/>
                        </div>
                    </div>
                    <div class="brand-card-body">
                        <div>
                            <div class="text-value text-danger">{{volumns.node.online}}</div>
                            <div class=" text-uppercase text-muted">Online</div>
                        </div>
                        <div>
                            <div class="text-value text-primary">{{volumns.node.all}}</div>
                            <div class="text-uppercase text-muted ">All</div>
                        </div>
                    </div>
                </div>
            </b-col>
            <b-col md="3" sm="6">
                <div class="brand-card">
                    <div class="brand-card-header bg-twitter">
                        <i class="fa fa-desktop"> Total Process</i>
                        <div class="chart-wrapper">
                            <social-box-chart :data="[1, 13, 9, 17, 34, 41, 38]"/>
                        </div>
                    </div>
                    <div class="brand-card-body">
                        <div>
                            <div class="text-value">{{volumns.process.started}}</div>
                            <div class="text-uppercase text-muted ">Started</div>
                        </div>
                        <div>
                            <div class="text-value">{{volumns.process.all}}</div>
                            <div class="text-uppercase text-muted ">All</div>
                        </div>
                    </div>
                </div>
            </b-col>
            <b-col md="3" sm="6">
                <div class="brand-card">
                    <div class="brand-card-header bg-linkedin">
                        <i class="fa fa-area-chart"></i>
                        <div class="chart-wrapper">
                            <social-box-chart :data="[78, 81, 80, 45, 34, 12, 40]"/>
                        </div>
                    </div>
                    <div class="brand-card-body">
                        <div>
                            <div class="text-value">{{volumns.info.CPU}}%</div>
                            <div class="text-uppercase text-muted">CPU</div>
                        </div>
                        <div>
                            <div class="text-value">{{ramShow(volumns.info.RAM)}}</div>
                            <div class="text-uppercase text-muted">RAM</div>
                        </div>
                    </div>
                </div>
            </b-col>
        </b-row>
    </div>
</template>

<script>
    import SocialBoxChart from "./SocialBoxChart";
    export default {
        name: 'dashboard',
        components: {SocialBoxChart},
        data: () => ({
            volumns: {
                info: {
                    CPU: "?",
                    RAM: "?"
                },
                node: {
                    all: "?",
                    online: "?"
                },
                process: {
                    all: "?",
                    started: '？'
                },
                version: {
                    "tag_name": "v2.0.0",
                    name: "",
                    "html_url": ""
                }
            },
        }),
        mounted() {
            this.queryDashboard()
        },
        methods: {
            queryDashboard() {
                let self = this;
                self.request("Dashboard加载中。", self.$axios.get("/admin/dashboard").then((res) => {
                    self.volumns = res;
                }).catch(e => {
                    self.$toast.error(e.message);
                }));
            }
        }
    }
</script>
