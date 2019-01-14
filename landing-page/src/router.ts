import Vue from 'vue';
import Router from 'vue-router';
import Home from './views/Home.vue';

Vue.use(Router);

export default new Router({
  mode: 'history',
  base: process.env.BASE_URL,
  scrollBehavior(to, from, savedPosition) {
    return { x: 0, y: 0 };
  },
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home,
      meta: {
        title: 'Home - Dynauth',
        description: 'Home page of Dynauth landing page',
      },
    },
    {
      path: '/dynamic-authentication',
      name: 'dynamic-authentication',
      component: () => import('./views/DynamicAuthentication.vue'),
      meta: {
        title: 'Dynamic Authentication - Dynauth',
        description: 'Outline of dynamic authentication by Dynauth',
      },
    },
    {
      path: '/mvp',
      name: 'mvp',
      component: () => import('./views/MVP.vue'),
      alias: '/password-manager',
      meta: {
        title: 'MVP - Dynauth',
        description: 'MVP description of Dynauth password manager',
      },
    },
    {
      path: '/password-manager',
      name: 'password-manager',
      component: () => import('./views/PasswordManager.vue'),
      meta: {
        title: 'Password Manager - Dynauth',
        description: 'Description of Dynauth Password Manager',
      },
    },
    // {
    //   path: '/about',
    //   name: 'about',
    //   component: () => import('./views/About.vue'),
    //   meta: {
    //     title: 'About - Dynauth',
    //     description: 'About page of Dynauth landing page',
    //   },
    // },
    {
      path: '/contact',
      name: 'contact',
      component: () => import('./views/Contact.vue'),
      meta: {
        title: 'Contact Me - Dynauth',
        description: 'Contact me if you: want to work with me, have input, have criticisms, or want a risky investment opportunity',
      },
    },
  ],
});
