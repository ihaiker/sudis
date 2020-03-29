import Vue from "vue"
import axios from 'axios'
import main from '../main'

axios.defaults.baseURL = process.env.VUE_APP_URL;
axios.defaults.timeout = 15000;
axios.defaults.withCredentials = true;
axios.defaults.headers.post['Content-Type'] = 'application/json;charset=UTF-8';

// 添加请求拦截器
axios.interceptors.request.use(function (config) {
    let token = localStorage.getItem('token');
    if (token) {
        config.headers['Authorization'] = token;
    }
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
            localStorage.removeItem("token");
            if (err.response.data && err.response.data.redirect) {
                window.location.href = err.response.data.redirect;
            } else {
                main.$router.push({path: '/signin', replace: true});
            }
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
