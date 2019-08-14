<template>
  <div>
    <Header />

    <Messages class="pt-4" />

    <b-container class="pt-4 pb-4" fluid>
      <b-row class="pb-3">
        <b-col>
          <b-button-toolbar>
            <b-button-group>
              <b-button variant="primary" @click.prevent="refreshPlayers">
                <font-awesome-icon icon="redo"></font-awesome-icon>
              </b-button>
            </b-button-group>
          </b-button-toolbar>
        </b-col>
      </b-row>

      <b-row>
        <b-col>
          <b-spinner
            v-show="refreshingPlayers"
          ></b-spinner>

          <b-table
            small
            striped
            hover
            class="players"
            :items="players"
            v-show="!refreshingPlayers"
            show-empty
          >
          </b-table>
        </b-col>
      </b-row>
    </b-container>
  </div>
</template>

<script>
import Vue from 'vue'
import { mapState, mapActions } from 'vuex'
import BootstrapVue from 'bootstrap-vue'
import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'
import { library } from '@fortawesome/fontawesome-svg-core'
import { faRedo } from '@fortawesome/free-solid-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import Header from '@/components/Core/Header'
import Messages from '@/components/Messages'

Vue.use(BootstrapVue)

library.add(faRedo)

export default {
  mounted() {
    this.refreshGame()

    this.refreshPlayers()
  },
  components: {
    FontAwesomeIcon, Header, Messages
  },
  computed: {
    ...mapState([
      'game',
      'refreshingPlayers', 'players'
    ])
  },
  methods: {
    ...mapActions([
      'refreshGame',
      'refreshPlayers'
    ])
  }
}
</script>

<style scoped>
.players {
  text-align: left;
}
</style>
