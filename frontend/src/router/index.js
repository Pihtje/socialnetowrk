import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import SignUpView from '../views/SignUpView.vue'
import FeedView from '../views/FeedView.vue'
import ProfileView from '../views/ProfileView.vue'
import CreatePostView from '../views/CreatePostView.vue'
import GroupsView from '../views/GroupsView.vue'
import SingelGroupView from '../views/SingleGroupView.vue'
import CreateGroupEventView from '@/views/CreateGroupEventView.vue'
import SingleEventView from '@/views/SingleEventView.vue'
const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/signup',
      name: 'signup',
      component: SignUpView
    },
    {
      path: '/feed',
      name: 'feed',
      component: FeedView
    },
    {
      path: '/profile/:userId',
      name: 'profile',
      component: ProfileView
    },
    
    {
      path: '/createPost',
      name: 'createPost',
      component: CreatePostView
    },
    {
      path: '/Groups',
      name: 'Groups',
      component: GroupsView
    },
    {
      path: '/SingleGroup/:groupId',
      name: 'SingleGroup',
      component: SingelGroupView
    },
    {
      path: '/createGroupEvent',
      name: 'CreateGroupEventView',
      component: CreateGroupEventView
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'login',
      component: LoginView
    },
  ]
})

export default router
