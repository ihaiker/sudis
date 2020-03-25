import 'core-js/es6/promise'
import 'core-js/es6/string'
import 'core-js/es7/array'

// import Vue from 'vue'
import BootstrapVue from 'bootstrap-vue'
import App from './App'
import router from './router'
import mixins from "./tools/mixins"
//全局插件组件
import plugins from './tools/plugins'
import http from "./tools/http"

Vue.use(BootstrapVue);
Vue.use(plugins);
Vue.mixin(mixins);

Vue.config.productionTip = false;

Vue.prototype.$axios = http.axios;
Vue.prototype.$form = http.form;


import {Alert, Confirm} from 'vue-m-dialog';

//https://mengdu.github.io/m-dialog/example/index.html#example
Alert.config({
    "title": "提示？", "show-close": false,
    "confromButtonText": "确定",
    "confirmButtonClassName": "btn btn-info btn-block"
});
Vue.prototype.$alert = Alert;

Confirm.config({
    "title": "确定？", "show-close": false,
    closeOnClickModal: false, closeOnPressEscape: false, showClose: false,
    cancelButtonText: "取消", cancelButtonClassName: "btn btn-light w-25",
    confromButtonText: "确定", confirmButtonClassName: "btn btn-info w-25",
});
Vue.prototype.$confirm = Confirm;

//https://github.com/bajian/vue-toast
import Toast from 'vue-bajiantoast'
import '@/assets/toast.css';

Vue.prototype.$toast = Toast;
Toast.config({
    duration: 3000,
    position: 'top right', showCloseBtn: true,
});



import VueParticles from 'vue-particles'
Vue.use(VueParticles);

let vm = new Vue({
    el: '#app',
    router,
    data: {loadingShow: false, loadingTitle: ""},
    template: '<App/>',
    components: {App}
});

export default vm
