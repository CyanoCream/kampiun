import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import './style.css'
import App from './App.vue'

import Landing from './views/Landing.vue'
import LiveScore from './views/LiveScore.vue'
import Bracket from './views/Bracket.vue'
import Dashboard from './views/Dashboard.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', name: 'landing', component: Landing },
    { path: '/live/:id', name: 'live', component: LiveScore, props: true },
    { path: '/bracket/:id', name: 'bracket', component: Bracket, props: true },
    { path: '/dashboard', name: 'dashboard', component: Dashboard },
  ],
})

createApp(App).use(router).mount('#app')
