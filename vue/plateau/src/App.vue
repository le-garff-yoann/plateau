<template>
  <div id="app">
    <router-view v-if="sessionCookie" />
    <Login v-else />
  </div>
</template>

<script>
import { mapState, mapMutations } from 'vuex'
import Login from '@/components/Login'
import consts from '@/constants/index.js'

export default {
  name: 'app',
  mounted() {
    this.setSessionCookie(window.$cookies.get(consts.serverName))
  },
  components: { Login },
  computed: {
    ...mapState([ 
      'sessionCookie'
    ])
  },
  methods: {
    ...mapMutations([
      'setSessionCookie'
    ])
  }
}
</script>

<style>
#app {
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
}
</style>
