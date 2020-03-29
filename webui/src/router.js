// import Vue from 'vue'
// import VueRouter from 'vue-router'

const AppContainer = () => import('@/containers/AppContainer');
const Login = () => import('@/views/Login');

const Dashboard = () => import('@/views/dashboard/Dashboard');
const Nodes = () => import('@/views/nodes/Nodes');
const Programs = () => import('@/views/program/ProgramList');
const ProgramInfo = () => import('@/views/program/ProgramInfo');
const ProgramLogger = () => import('@/views/program/ProgramLogger');
const Tags = () => import('@/views/tags/Tags');
const Users = () => import('@/views/users/Users');
const Systems = () => import('@/views/system/Systems');
const Logs = () => import('@/views/logs/Logs');

Vue.use(VueRouter);

export default new VueRouter({
    mode: 'hash',
    linkActiveClass: 'open active',
    scrollBehavior: () => ({y: 0}),
    routes: [
        {path: "/signin", name: 'signin', component: Login},
        {
            path: "/admin", component: AppContainer,
            children: [
                {path: "", redirect: "dashboard"},

                {path: "dashboard", component: Dashboard},
                {path: "nodes", component: Nodes},
                {path: "programs", component: Programs},
                {path: "programs/:tag", component: Programs},

                {path: "program", component: ProgramInfo},
                {path: "program/logs", component: ProgramLogger},

                {path: "tags", component: Tags},
                {path: "users", component: Users},
                {path: "systems", component: Systems},
                {path: "logs", component: Logs},


                {path: '*', redirect: 'dashboard'}
            ]
        },
        {path: '*', redirect: '/admin'}
    ]
})
