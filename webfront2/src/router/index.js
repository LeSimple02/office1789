import { createRouter, createWebHistory, createWebHashHistory } from 'vue-router'
import { gls } from '@/stores/global'
import HomeView from '@/views/HomeView.vue'
import MailView from '@/views/MailView.vue'


const useHash = typeof window !== 'undefined' && (window.location.protocol === 'file:' || window.cordova)
const router = createRouter({
  history: useHash ? createWebHashHistory(import.meta.env.BASE_URL) : createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/LoginView.vue')
    },
    {
      path: '/createaccount',
      name: 'create',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('../views/CreateAccount.vue')
    },
    {
      path: '/forgot',
      name: 'forgot',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('../views/ForgotP.vue')
    },
    {
      path: '/reset-password',
      name: 'reset-password',
      component: () => import('../views/ResetPassword.vue')
    },
    {
      path: '/legalesmentions',
      name: 'Légales',
      component: () => import('../views/MentionsView.vue')
    },
    {
      path: '/mail',
      name: 'mail',
      component: MailView,
      children: [
      {
        path: 'draft',
        component: MailView,
      },{
        path: 'send',
        component: MailView,
      }, {
        path: 'trash',
        component: MailView,
      }]
    },
    {
      path: '/drive',
      name: 'drive',
      component: () => import('../views/DriveView.vue'),
       children: [
      {
        path: 'trash',
        component: () => import('../views/DriveView.vue')
      }]
    },
    {
      path: '/chat',
       children: [
      {
	path: '',
	name: 'chat',
	component: () => import('../views/ChatView.vue')
      },
      {
	path: '/:user',
	component: () => import('../views/ChatView.vue'),
      }]
    },
    {
      path: '/account',
      name: 'account',
      component: () => import('../views/AccountView.vue')
    },
    {
      path: '/account/edit',
      name: 'account-edit',
      component: () => import('../views/AccountView.vue')
    },
    {
      path: '/account/subscription-success',
      name: 'subscription-success',
      component: () => import('../components/SubscriptionSuccess.vue')
    },
    {
      path: '/account/organization',
      name: 'organization',
      component: () => import('../components/OrganizationPanel.vue')
    },
    {
      path: '/account/custom-domain',
      name: 'custom-domain',
      component: () => import('../components/CustomDomainPanel.vue')
    },
    {
      path: '/about',
      name: 'about',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('../views/AboutView.vue')
    },
    {
	path: '/contact',
	component: () => import('../views/ContactView.vue')
    },
     {
	path: '/calendar',
	component: () => import('../views/CalendarView.vue')
    },
    {
      path: '/admin',
      name: 'admin',
      component: () => import('../views/AdminPanel.vue')
    }
  ]
})

const isDevMode = typeof window !== 'undefined' && new URLSearchParams(window.location.search).get('dev') === '1'

const protectedPaths = ['/mail', '/drive', '/chat', '/calendar', '/account', '/admin']

router.beforeEach((to) => {
  const store = gls()
  const needsAuth = protectedPaths.some(p => to.path === p || to.path.startsWith(p + '/'))
  if (needsAuth && store.log != 1 && !isDevMode) {
    return '/login'
  }
  if (to.path === '/login' && store.log == 1) {
    return '/mail'
  }
})

export default router
