<template>
  <div>
    <Header />

    <Messages class="pt-4" />

    <b-container class="pt-4 pb-4" fluid>
      <b-row class="pb-3">
        <b-col>
          <b-button-toolbar>
            <b-button-group>
              <b-button @click.prevent="refreshMatchs" variant="primary">
                <font-awesome-icon icon="redo"></font-awesome-icon>
              </b-button>

              <b-button variant="success" v-b-modal.new-match>New</b-button>
            </b-button-group>
          </b-button-toolbar>
        </b-col>
      </b-row>

      <b-row>
        <b-col>
          <b-spinner
            v-show="refreshingMatchs"
          ></b-spinner>

          <b-table
            class="matchs"
            small
            striped
            hover
            v-show="!refreshingMatchs"
            :items="matchs"
          >
            <template slot="id" slot-scope="data">
              <b-link v-bind:to="`/match/${data.value}`">{{ data.value }}</b-link>
            </template>
          </b-table>
        </b-col>
      </b-row>
    </b-container>

    <b-modal id="new-match" title="New match" @ok="doCreateMatch" size="xl">
      <label for="input-live">Number of players required: {{ newMatchNumberOfPlayersRequired }}</label>

      <b-form-input
        v-model="newMatchNumberOfPlayersRequired"
        :run="!newMatchNumberOfPlayersRequired ? newMatchNumberOfPlayersRequired = game.min_players : true"
        type="range"
        v-bind:min="game.min_players"
        v-bind:max="game.max_players"
      ></b-form-input>
    </b-modal>
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

    this.refreshMatchs()
  },
  components: {
    FontAwesomeIcon, Header, Messages
  },
  data() {
    return {
      newMatchNumberOfPlayersRequired: 0
    }
  },
  computed: {
    ...mapState([
      'game',
      'refreshingMatchs', 'matchs'
    ])
  },
  methods: {
    ...mapActions([
      'refreshGame',
      'refreshMatchs', 'createMatch'
    ]),
    doCreateMatch() {
      this.createMatch(this.newMatchNumberOfPlayersRequired)
    }
  }
}
</script>

<style scoped>
.matchs {
  text-align: left;
}
</style>
