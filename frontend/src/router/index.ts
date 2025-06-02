import { createRouter, createWebHistory } from 'vue-router'
import AuthView from '@/views/AuthView.vue'
import PanelView from '@/views/PanelView.vue'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', component: AuthView },
    {
      path: '/panel',
      name: 'panel',
      component: PanelView
    },
    {
      path: '/panel/profile',
      name: 'UserProfileView',
      component: () => import('@/views/UserProfileView.vue')
    }
  ]
})

router.beforeEach((to, from, next) => {
  const store = useAuthStore()
  const publicPaths = ['/', '/register', '/login']
  const authRequired = !publicPaths.includes(to.path)

  if (authRequired && !store.token) {
    return next('/')
  }

  next()
})

export default router
