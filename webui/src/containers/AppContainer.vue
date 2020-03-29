<template>
    <div class="app">
        <Loadings/>
        <AppHeader fixed>
            <SidebarToggler class="d-lg-none" display="md" mobile/>
            <b-link class="navbar-brand" to="/admin">
                <img class="navbar-brand-full" src="@/assets/images/logo.png" width="89" height="25" alt="Sudis Logo">
                <img class="navbar-brand-minimized" src="/favicon.ico" width="30" height="30" alt="Sudis Logo">
            </b-link>

            <SidebarToggler class="d-md-down-none" display="lg"/>

            <b-navbar-nav class="d-md-down-none">
                <b-nav-item class="px-3" to="/admin">分布式守护进程管理器</b-nav-item>
            </b-navbar-nav>

            <b-navbar-nav class="ml-auto d-md-down-none">
                <b-nav-item class="px-3" href="http://sudis.renzhen.la" target="_blank">/官方主页/</b-nav-item>
                <b-nav-item class="px-3" href="https://github.com/ihaiker/sudis/wiki" target="_blank">/说明文档/</b-nav-item>
                <b-nav-item class="px-3" href="https://github.com/ihaiker/sudis/issues/new" target="_blank">/提交BUG/</b-nav-item>
                <b-nav-item class="px-3" href="https://github.com/ihaiker/sudis" target="_blank">/源码下载/</b-nav-item>
                <b-nav-item class="px-3" href="http://shui.renzhen.la" target="_blank">/关于作者/</b-nav-item>

                <AppHeaderDropdown right no-caret class="mr-3">
                    <template slot="header">
                        <img src="@/assets/images/logo2.png" class="img-avatar"/>
                    </template>
                    <template slot="dropdown">
                        <b-dropdown-item><i class="fa fa-user-circle"/> <strong>{{userName}}</strong></b-dropdown-item>
                        <!--<b-dropdown-item><i class="fa fa-shield"/> 修改密码</b-dropdown-item> -->
                        <b-dropdown-item @click="logout"><i class="fa fa-lock"/> 退出登录</b-dropdown-item>
                    </template>
                </AppHeaderDropdown>

                <!--
                <b-nav-item>
                    <i class="icon-bell"/>
                    <b-badge pill variant="danger">5</b-badge>
                </b-nav-item>
                -->
            </b-navbar-nav>
            <!--
            <AsideToggler class="d-none d-lg-block"/>
            <AsideToggler class="d-lg-none" mobile/>
            -->
        </AppHeader>
        <div class="app-body">
            <Sidebar fixed>
                <SidebarHeader/>
                <SidebarForm/>
                <SidebarNav :navItems="nav"/>
                <SidebarFooter/>
                <SidebarMinimizer/>
            </Sidebar>
            <main class="main">
                <router-view/>
            </main>
            <Aside fixed>
                <AppAside/>
            </Aside>
        </div>
        <AppFooter>
            <div></div>
            <div class="ml-auto">
                <a href="http://sudis.renzhen.la" target="_blank">Sudis</a>
                <span class="ml-1">&copy; 2019.</span>
            </div>
            <div class="ml-auto">
                <span class="mr-1">Powered by</span>
                <a href="http://shui.renzhen.la" target="_blank">Haiker</a>
            </div>
        </AppFooter>
    </div>
</template>

<script>
    import {Header as AppHeader, Footer as AppFooter} from '@coreui/vue'
    import {Sidebar, SidebarFooter, SidebarForm, SidebarHeader, SidebarMinimizer, SidebarNav, SidebarToggler} from '@coreui/vue'
    import {Aside, AsideToggler} from '@coreui/vue'
    import {HeaderDropdown as AppHeaderDropdown} from '@coreui/vue'

    import AppAside from './AppAside'
    import Loadings from "../plugins/loadings";

    export default {
        name: 'AppContainer',
        components: {
            Loadings,
            AppHeader, AppFooter,
            Aside, AsideToggler, AppAside,
            Sidebar, SidebarForm, SidebarFooter, SidebarToggler,
            SidebarHeader, SidebarNav, SidebarMinimizer, AppHeaderDropdown
        },
        mounted() {
            //this.queryTags();
        },
        methods: {
            queryTags() {
                this.$axios.get("/admin/tag/list")
                    .then(this.appendNavs);
            },
            appendNavs(tags) {
                for (let idx in tags) {
                    let tag = tags[idx];
                    this.nav[2].children.push({
                        name: '标签：' + tag.name,
                        url: '/admin/programs/' + tag.name,
                    });
                }
            },
            logout() {
                localStorage.removeItem("token");
                this.$router.push({path: "/login"});
            }
        },
        computed: {
            userName() {
                return localStorage.getItem("x-user");
            }
        },
        data: () => ({
            nav: [
                {
                    name: 'Dashboard',
                    url: '/admin/dashboard',
                    icon: 'icon-graph',
                    badge: {
                        variant: 'primary',
                        text: '新'
                    }
                },
                {
                    name: '节点管理',
                    url: '/admin/nodes',
                    icon: 'icon-puzzle'
                },
                {
                    name: '进程管理',
                    icon: 'fa fa-microchip',
                    url: '/admin/programs',
                },
                {
                    name: '标签管理',
                    url: '/admin/tags',
                    icon: 'fa fa-tags'
                },
                {
                    name: '用户管理',
                    url: '/admin/users',
                    icon: 'fa fa-user'
                },
                {
                    name: '系统管理',
                    url: '/admin/systems',
                    icon: 'icon-settings'
                },
                /*{
                    name: '查阅日志',
                    url: '/admin/logs',
                    icon: 'icon-list'
                }*/
            ]
        })
    }
</script>
