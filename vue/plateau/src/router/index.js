import Vue from 'vue'
import Router from 'vue-router'
import Matchs from '@/components/Core/Matchs'
import Match from '@/components/Core/Match'
import Players from '@/components/Core/Players'
import About from '@/components/Core/About'
import NotFound from '@/components/NotFound'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/matchs',
      alias: '/',
      name: 'Matchs',
      component: Matchs
    },
    {
      path: '/match/:id',
      component: Match
    },
    {
      path: '/players',
      name: 'Players',
      component: Players
    },
    {
      path: '/about',
      name: 'About',
      component: About
    },
    {
      path: '*',
      component: NotFound
    }
  ]
})
