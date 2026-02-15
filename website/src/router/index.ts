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
      path: '/report_card/:name',
      name: 'report_card',
      component: () => import('../views/report_card.vue')
    },
    {
      path: '/compare',
      name: 'compare',
      component: () => import('../views/compare.vue')
    },
    {
      path: '/coming-soon',
      name: 'coming-soon',
      component: () => import('../views/coming_soon.vue')
    },
    {
      path: '/features',
      name: 'features',
      component: () => import('../views/features.vue')
    },
    {
      path: '/matrix',
      name: 'matrix',
      component: () => import('../views/matrix.vue')
    }
  ],
})

export default router
