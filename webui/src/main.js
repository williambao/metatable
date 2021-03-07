// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'
import ViewUI from 'view-design';
// import 'iview/dist/styles/iview.css';
import '../my-theme/index.less';

import Moment from 'vue-moment';
import axios from './axios'
import store from './store/index.js'
import "babel-polyfill";

Vue.config.productionTip = false

Vue.use(ViewUI);
Vue.use(Moment);

Vue.prototype.$config = process.env;

/* eslint-disable no-new */
// new Vue({
//   el: '#app',
//   router,
//   template: '<App/>',
//   components: { App }
// })

new Vue({
    axios,
    store,
    router,
    render: h => h(App)
}).$mount('#app');