import Vue from 'vue';
import VueRouter from 'vue-router';
import App from './App.vue';
import Home from './Home.vue';
import SearchForm from './components/SearchForm.vue';
import About from './components/About.vue';
import Results from './components/Results.vue';

Vue.use(VueRouter);

Vue.config.productionTip = false;

const router = new VueRouter({
  routes: [
    { path: '/', component: Home },
    { path: '/search', component: SearchForm },
    { path: '/about', component: About },
    { path: '/results/:ID', component: Results }
  ]
});

new Vue({
  router,
  render: h => h(App),
}).$mount('#app');






