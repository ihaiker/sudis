import Vue from "vue"
import axios from 'axios'
import main from '../main'

axios.defaults.baseURL = process.env.VUE_APP_URL;
axios.defaults.timeout = 15000;
axios.defaults.withCredentials = true;
axios.defaults.headers.post['Content-Type'] = 'application/json;charse=UTF-8';

// 添加请求拦截器
axios.interceptors.request.use(function (config) {
    config.headers['x-ticket'] = localStorage.getItem('x-ticket');
    config.headers['x-user'] = localStorage.getItem('x-user');
    return config
});

//添加响应拦截器
axios.interceptors.response.use((response) => {
    if (response.status === 200 && response.data) {
        return response.data
    }
    return response
}, function (err) {
    if (err.response) {
        if (err.response.status === 401) {
            localStorage.removeItem("x-ticket");
            localStorage.removeItem("x-user");
            main.$router.push({path: '/login', replace: true});
        } else {
            return Promise.reject(err.response.data)
        }
    } else {
        return Promise.reject({e: err, message: err.message})
    }
});

let config = {
    transformRequest: [function (data) {
        let ret = '';
        for (let it in data) {
            ret += encodeURIComponent(it) + '=' + encodeURIComponent(data[it]) + '&'
        }
        return ret
    }],
    headers: {
        'Content-Type': 'application/x-www-form-urlencoded'
    }
};

export default {
    axios: axios,
    form: config
}
