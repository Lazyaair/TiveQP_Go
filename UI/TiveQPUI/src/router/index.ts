import { createRouter, createWebHistory } from 'vue-router'
import Login from '../components/Login.vue'
import RoleSelect from '../components/RoleSelect.vue'
import UserDashboard from '../views/UserDashboard.vue'
import OwnerDashboard from '../views/OwnerDashboard.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      redirect: '/login'
    },
    {
      path: '/login',
      name: 'login',
      component: Login
    },
    {
      path: '/role-select',
      name: 'role-select',
      component: RoleSelect
    },
    {
      path: '/user-dashboard',
      name: 'user-dashboard',
      component: UserDashboard,
      meta: { requiresAuth: true, role: 'user' }
    },
    {
      path: '/owner-dashboard',
      name: 'owner-dashboard',
      component: OwnerDashboard,
      meta: { requiresAuth: true, role: 'owner' }
    }
  ]
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const isAuthenticated = localStorage.getItem('token') // 检查是否已登录
  const userRole = localStorage.getItem('role') // 获取用户角色

  if (to.meta.requiresAuth && !isAuthenticated) {
    // 需要认证但未登录，重定向到登录页
    next({ name: 'login' })
  } else if (to.meta.role && to.meta.role !== userRole) {
    // 角色不匹配，重定向到角色选择页
    next({ name: 'role-select' })
  } else {
    next()
  }
})

export default router 