<template>
  <div>
    <b-navbar fixed="top" type="dark" variant="dark" sticky>
      <b-navbar-brand class="font-weight-light">
        <em>{{ game.name }}</em>
      </b-navbar-brand>

      <b-collapse id="nav-collapse" is-nav>
        <b-navbar-nav>
          <b-nav-item>
            <b-link to="/">Matchs</b-link>
          </b-nav-item>

          <b-nav-item>
            <b-link to="/players">Players</b-link>
          </b-nav-item>

          <b-nav-item>
            <b-link to="/about">About</b-link>
          </b-nav-item>
        </b-navbar-nav>

        <b-navbar-nav class="ml-auto">
          <b-nav-item-dropdown v-bind:text="`@${username}`" right>
            <b-dropdown-item @click.prevent="logout">Sign out</b-dropdown-item>
          </b-nav-item-dropdown>
        </b-navbar-nav>
      </b-collapse>
    </b-navbar>
  </div>
</template>

<script>
import Vue from 'vue'
import { mapState, mapMutations, mapActions } from 'vuex'
import BootstrapVue from 'bootstrap-vue'
import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

Vue.use(BootstrapVue)

export default {
  created() {
    this.setGlobalMessage(null)
    this.setGlobalError(null)
  },
  mounted() {
    this.refreshGame()
  },
  computed: {
    ...mapState([
      'localCookie',
      'game'
    ]),
    username() {
      return this.localCookie != null && 'username' in this.localCookie
        ? this.localCookie.username
        : null
    }
  },
  methods: {
    ...mapMutations([
      'setGlobalMessage', 'setGlobalError'
    ]),
    ...mapActions([
      'refreshGame', 'logout'
    ])
  }
}
</script>

<style scoped>
.router-link-exact-active {
  font-weight: bolder;
}
</style>
