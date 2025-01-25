import { createRouter, createWebHistory } from 'vue-router'
import Login from '../components/Login.vue'
import RoleSelect from '../components/RoleSelect.vue'
import ShopMap from '../components/ShopMap.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'login',
      component: Login
    },
    {
      path: '/role-select',
      name: 'role-select',
      component: RoleSelect,
      beforeEnter: (to, from, next) => {
        const token = localStorage.getItem('token')
        if (!token) {
          next('/')
        } else {
          next()
        }
      }
    },
    {
      path: '/map',
      name: 'map',
      component: ShopMap,
      beforeEnter: (to, from, next) => {
        const token = localStorage.getItem('token')
        const role = localStorage.getItem('userRole')
        if (!token) {
          next('/')
        } else if (!role) {
          next('/role-select')
        } else {
          next()
        }
      }
    }
  ]
})

// 全局前置守卫
router.beforeEach((to, from, next) => {
  console.log('Navigation to:', to.path)
  console.log('Token:', localStorage.getItem('token'))
  console.log('Role:', localStorage.getItem('userRole'))
  next()
})

export default router 