import Vue from 'vue';
import VueRouter from 'vue-router';
import { router } from '@router.js';

import api from './api.js';

Vue.config.productionTip = false;

// aliasing
Vue.prototype.$http = api;

new Vue({
  router,
  render: h => h(App),
}).$mount('#app');
