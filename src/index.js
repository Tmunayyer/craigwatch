import Vue from 'vue';
import VueRouter from 'vue-router';
import App from './App.vue';
import Home from './Home.vue';
import Results from './Results.vue';

Vue.use(VueRouter);

Vue.config.productionTip = false;

const router = new VueRouter({
  routes: [
    { path: '/', component: Home },
    { path: '/result/:ID', component: Results }
  ]
});

new Vue({
  router,
  render: h => h(App),
}).$mount('#app');
