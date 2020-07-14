import Vue from 'vue';
import VueRouter from 'vue-router';
import App from './App.vue';
import Home from './Home.vue';
import Results from './Results.vue';

import api from './api.js';

Vue.use(VueRouter);

Vue.config.productionTip = false;

const router = new VueRouter({
  routes: [
    { path: '/', component: Home },
    { path: '/result/:ID', component: Results }
  ]
});

// aliasing
Vue.prototype.$http = api;

new Vue({
  router,
  render: h => h(App),
}).$mount('#app');
