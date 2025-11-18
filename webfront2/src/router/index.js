import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '@/views/HomeView.vue'
import MailView from '@/views/MailView.vue'


const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
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
      children: [
      {
	path: '',
	name: 'account',
	component: ()=>import('../views/AccountView.vue')
      },
      {
	path: 'edit',
	component: ()=>import('../views/AccountView.vue'),
      },
      {
	path: 'subscription-success',
	name: 'subscription-success',
	component: ()=>import('../components/SubscriptionSuccess.vue'),
      }]
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
    }
  ]
})

export default router
