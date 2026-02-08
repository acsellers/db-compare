import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'landing',
      component: () => import('../views/landing.vue')
    },
    {
      path: '/benchmarks',
      name: 'benchmarks',
      component: () => import('../views/benchmarks.vue')
    },
    {
      path: '/libraries',
      name: 'libraries',
      component: () => import('../views/libraries.vue')
    },
    {
      path: '/examples',
      name: 'examples',
      component: () => import('../views/examples.vue')
    },
    {
      path: '/example/:name',
      name: 'example',
      component: () => import('../views/example.vue')
    },
    {
      path: '/compare',
      name: 'compare',
      component: () => import('../views/compare.vue')
    }
  ],
})

export default router
