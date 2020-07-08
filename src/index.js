import Vue from 'vue'
import VueRouter from 'vue-router'
import App from './App.vue'
import CreateSearch from './components/CreateSearch.vue'
import About from './components/About.vue'
import Results from './components/Results.vue'


Vue.use(VueRouter)


Vue.config.productionTip = false


const router = new VueRouter({
  routes: [
    { path: '/search', component: CreateSearch },
    { path: '/about', component: About },
    { path: '/results', component: Results }
  ]
})




new Vue({
  router,
  render: h => h(App),
}).$mount('#app')






