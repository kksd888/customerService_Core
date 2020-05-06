import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'

import '../src/assets/public.css'

// 主题
import ElementUI from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css'
Vue.use(ElementUI);

// cookie
import VueCookies from 'vue-cookies'
Vue.use(VueCookies);

// httpclient
import axios from './plugin/axios'
import VueAxios from 'vue-axios'
Vue.use(VueAxios, axios);

Vue.config.productionTip = false;

new Vue({
  router,
  store,
  render: function (h) { return h(App) }
}).$mount('#app');
