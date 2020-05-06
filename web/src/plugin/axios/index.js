import axios from 'axios'
import {Notification} from 'element-ui'

const baseHost = "localhost:5000";
const apiHost = "http://" + baseHost;
const wsHost = "ws://" + baseHost;

// 创建一个 axios 实例
const service = axios.create({
    baseURL: apiHost,
    timeout: 30000
});

service.apiHost = apiHost;
service.wsHost = wsHost;

// 请求拦截器
service.interceptors.request.use(config => {
    if (window.$cookies.isKey("token")) {
        config.headers['Authentication'] = window.$cookies.get("token");
    }
    return config;
}, error => {
    return Promise.reject(err);
});

// 响应拦截器
service.interceptors.response.use(
    response => {
        return response;
    },
    error => {
        if (error && error.response) {
            switch (error.response.status) {
                case 400:
                    error.message = '请求错误';
                    break;
                case 401:
                    error.message = '未授权，请登录';
                    break;
                case 403:
                    error.message = '拒绝访问';
                    break;
                case 404:
                    error.message = `请求地址出错: ${error.response.config.url}`;
                    break;
                case 408:
                    error.message = '请求超时';
                    break;
                case 500:
                    error.message = '服务器内部错误';
                    break;
                case 501:
                    error.message = '服务未实现';
                    break;
                case 502:
                    error.message = '网关错误';
                    break;
                case 503:
                    error.message = '服务不可用';
                    break;
                case 504:
                    error.message = '网关超时';
                    break;
                case 505:
                    error.message = 'HTTP版本不受支持';
                    break;
                default:
                    break;
            }
        }

        Notification.warning({
            title: '服务器接口异常',
            message: `${error}`
        });

        if (error.response.status === 401) {
            window.$cookies.remove('token');
            window.location.reload();
        }

        return Promise.reject(error)
    }
);

export default service