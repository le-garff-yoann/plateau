<template>
  <b-container>
    <b-row align-h="center">
      <b-col md="5" class="pt-4">
        <b-img src="@/assets/logo.png"></b-img>
      </b-col>
    </b-row>

    <Messages />

    <b-row align-h="center">
      <b-col md="5">
        <b-spinner
          v-show="loggingIn || registerIn"
        ></b-spinner>
        <b-form
          @submit.prevent="doAction"
          v-show="!loggingIn && !registerIn"
        >
          <b-form-group>
            <b-form-input
              v-model="username"
              type="text"
              placeholder="Username"
              required
            ></b-form-input>
          </b-form-group>

          <b-form-group>
            <b-form-input
              v-model="password"
              type="password"
              placeholder="Password"
            ></b-form-input>
          </b-form-group>

          <b-button
            @click="action = 'in'"
            type="submit"
            variant="primary"
            block
          >Sign in</b-button>
          <b-button
            @click="action = 'up'"
            type="submit"
            variant="dark"
            block
          >Sign up</b-button>
        </b-form>
      </b-col>
    </b-row>
  </b-container>
</template>

<script>
import Vue from 'vue'
import { mapState, mapMutations, mapActions } from 'vuex'
import BootstrapVue from 'bootstrap-vue'
import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'
import Messages from '@/components/Messages'

Vue.use(BootstrapVue)

export default {
  created() {
    this.setGlobalMessage(null)
    this.setGlobalError(null)
  },
  components: { Messages },
  data() {
    return {
      action: null,
      username: null,
      password: null
    }
  },
  computed: {
    ...mapState([
      'globalMessage', 'globalError',
      'loggingIn', 'registerIn'
    ])
  },
  methods: {
    ...mapMutations([
      'setGlobalMessage', 'setGlobalError'
    ]),
    ...mapActions([
      'login', 'register'
    ]),
    doAction() {
      const userinfo = {
        username: this.username,
        password: this.password
      }

      this.action == 'in'
        ? this.login(userinfo)
        : this.register(userinfo)
    }
  }
}
</script>
