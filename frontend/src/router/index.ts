import { createRouter, createWebHashHistory} from 'vue-router'
import MainPage from '../components/Main/Index.vue'
import SecondPage from '../pages/Home/Index.vue'
import TipsPage from '../pages/Tips/Index.vue'

const routes = [
  { path: '/', component: MainPage },
  { path: '/second', component: SecondPage },
  { path: '/tips', component: TipsPage }
]

export default createRouter({
  history: createWebHashHistory(), // Use createWebHashHistory for hash-based routing
  routes
})