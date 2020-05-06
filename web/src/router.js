import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router);

const routers = new Router({
    mode: 'hash',
    base: process.env.BASE_URL,
    routes: [
        {
            path: '/',
            redirect: {name: 'login'},
        },
        {
            name: 'login',
            path: '/login',
            meta: {
                auth: false,
                title: "登录",
            },
            component: function () {
                return import('./views/Login.vue')
            }
        },
        {
            name: 'home',
            path: '/home',
            meta: {
                auth: true,
                title: "聊天面板",
            },
            component: function () {
                return import('./views/Home.vue')
            }
        }
    ]
});

routers.beforeEach((to, from, next) => {
    // 验证当前路由所有的匹配中是否需要有登录验证的
    if (to.matched.some(r => r.meta.auth)) {
        // cookie里是否存有token作为验证是否登录的条件
        if (window.$cookies.isKey("token")) {
            next()
        } else {
            // 没有登录的时候跳转到登录界面
            next({
                name: 'login',
                query: {
                    redirect: to.fullPath
                }
            });
        }
    } else {
        next()
    }
});

routers.afterEach(to => {
    if (to.meta.title) {
        document.title = "金色世纪在线客服 - " + to.meta.title;
    }
});


export default routers