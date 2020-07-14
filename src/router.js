import Vue from "vue";
import VueRouter from "vue-router";

import App from './App.vue';
import Home from './Home.vue';
import Results from './Results.vue';

Vue.use(VueRouter);

export const routes = [
    { path: '/', component: Home },
    { path: '/result/:ID', component: Results }
];

export const router = new VueRouter({
    routes
});