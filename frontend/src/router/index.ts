import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory('/admin/'),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/Login.vue'),
      meta: { public: true },
    },
    {
      // Builder: fullscreen route, no AdminLayout/sidebar.
      path: '/builder/:id',
      name: 'builder',
      component: () => import('@/views/Builder.vue'),
    },
    {
      path: '/',
      component: () => import('@/layouts/AdminLayout.vue'),
      children: [
        { path: '', name: 'dashboard', component: () => import('@/views/Dashboard.vue') },
        { path: 'posts', name: 'posts', component: () => import('@/views/Posts.vue') },
        { path: 'posts/new', name: 'post-new', component: () => import('@/views/PostEditor.vue') },
        { path: 'posts/:id/edit', name: 'post-edit', component: () => import('@/views/PostEditor.vue') },
        { path: 'projects', name: 'projects', component: () => import('@/views/Projects.vue') },
        { path: 'projects/new', name: 'project-new', component: () => import('@/views/ProjectEditor.vue') },
        { path: 'projects/:id/edit', name: 'project-edit', component: () => import('@/views/ProjectEditor.vue') },
        { path: 'open-source', name: 'open-source', component: () => import('@/views/OpenSource.vue') },
        { path: 'open-source/new', name: 'open-source-new', component: () => import('@/views/OpenSourceEditor.vue') },
        { path: 'open-source/:id/edit', name: 'open-source-edit', component: () => import('@/views/OpenSourceEditor.vue') },
        { path: 'products', name: 'products', component: () => import('@/views/Products.vue') },
        { path: 'products/new', name: 'product-new', component: () => import('@/views/ProductEditor.vue') },
        { path: 'products/:id/edit', name: 'product-edit', component: () => import('@/views/ProductEditor.vue') },
        { path: 'inbox', name: 'inbox', component: () => import('@/views/Inbox.vue') },
        { path: 'pages', name: 'pages', component: () => import('@/views/LandingPages.vue') },
        { path: 'users', name: 'users', component: () => import('@/views/Users.vue'), meta: { super: true } },
        { path: 'users/new', name: 'user-new', component: () => import('@/views/UserEditor.vue'), meta: { super: true } },
        { path: 'users/:id/edit', name: 'user-edit', component: () => import('@/views/UserEditor.vue'), meta: { super: true } },
        { path: 'profile', name: 'profile', component: () => import('@/views/Profile.vue') },
      ],
    },
  ],
})

router.beforeEach((to, _from, next) => {
  const token = localStorage.getItem('access_token')
  if (to.meta.public) {
    next()
  } else if (!token) {
    next('/login')
  } else if (to.meta.super && localStorage.getItem('is_super') !== '1') {
    next('/')
  } else {
    next()
  }
})

export default router
