import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory('/admin/'),
  routes: [
    {
      path: '/admin/login',
      name: 'login',
      component: () => import('@/views/Login.vue'),
      meta: { public: true },
    },
    {
      path: '/admin',
      component: () => import('@/layouts/AdminLayout.vue'),
      children: [
        { path: '', name: 'dashboard', component: () => import('@/views/Dashboard.vue') },
        { path: 'posts', name: 'posts', component: () => import('@/views/Posts.vue') },
        { path: 'posts/new', name: 'post-new', component: () => import('@/views/PostEditor.vue') },
        { path: 'posts/:id/edit', name: 'post-edit', component: () => import('@/views/PostEditor.vue') },
        { path: 'projects', name: 'projects', component: () => import('@/views/Projects.vue') },
        { path: 'open-source', name: 'open-source', component: () => import('@/views/OpenSource.vue') },
        { path: 'products', name: 'products', component: () => import('@/views/Products.vue') },
        { path: 'inbox', name: 'inbox', component: () => import('@/views/Inbox.vue') },
        { path: 'settings', name: 'settings', component: () => import('@/views/Settings.vue') },
      ],
    },
  ],
})

router.beforeEach((to, _from, next) => {
  const token = localStorage.getItem('access_token')
  if (to.meta.public) {
    next()
  } else if (!token) {
    next('/admin/login')
  } else {
    next()
  }
})

export default router
