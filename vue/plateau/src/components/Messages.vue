<template>
  <b-container v-show="globalMessage != null || globalError != null">
    <b-row align-h="center">
      <b-col md="5">
        <b-alert
          variant="success"
          dismissible
          fade
          @dismissed="setGlobalMessage(null)"
          :show="globalMessage != null"
        >
          <p class="mb-0">{{ globalMessage }}</p>
        </b-alert>

        <b-alert
          variant="warning"
          dismissible
          fade
          @dismissed="setGlobalError(null)"
          :show="globalError != null"
        >
          <p class="mb-0">{{ globalError }}</p>
        </b-alert>
      </b-col>
    </b-row>
  </b-container>
</template>

<script>
import Vue from 'vue'
import { mapState, mapMutations } from 'vuex'
import BootstrapVue from 'bootstrap-vue'
import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

Vue.use(BootstrapVue)

export default {
  created() {
    this.setGlobalMessage(null)
    this.setGlobalError(null)
  },
  computed: {
    ...mapState([
      'globalMessage', 'globalError'
    ])
  },
  methods: {
    ...mapMutations([
      'setGlobalMessage', 'setGlobalError'
    ])
  }
}
</script>
