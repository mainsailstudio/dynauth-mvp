// Basic site imports like Bootstrap and jQuery
import 'bootstrap';
import './assets/vendor/jquery/jquery.min.js';

// Vue imports
import Vue from 'vue';
import App from './App.vue';
import router from './router';
import './registerServiceWorker';

// Axios import
import axios from 'axios';

Vue.config.productionTip = false;

// Change the page title and meta based on the router attributes
router.beforeEach((to, from, next) => {
  document.title = to.meta.title;
  // document.getElementsByTagName('meta')['description'].content = to.meta.description;
  next();
});

new Vue({
  router,
  render: (h) => h(App),
}).$mount('#app');

// APP CONSTANTS
const axiosConfig = {
  baseURL: 'http://localhost:8080',
  timeout: 30000,
};
Vue.prototype.axios = axios.create(axiosConfig);

const testApiURL = 'http://localhost:8080';
const realApiURL = '';
