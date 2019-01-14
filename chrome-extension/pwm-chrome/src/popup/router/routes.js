import PageIndex from './pages/Index'
import About from './pages/About'
import DynamicAuthentication from './pages/Authenticate'

export default [
  {
    path: '/',
    component: PageIndex
  },
  {
    path: '/about',
    component: About
  },
  {
    path: '/dynauth',
    component: DynamicAuthentication
  }
]
