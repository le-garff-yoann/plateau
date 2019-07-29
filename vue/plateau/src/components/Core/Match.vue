<template>
  <div>
    <Header />

    <Messages class="pt-4" />

    <b-container class="pt-4 pb-4" fluid>
      <b-card
        bg-variant="light"
        class="text-left"
      >
        <h6 slot="header">
          Match <em>{{ id }}</em>
        </h6>

        <b-container fluid>
          <b-row class="pb-3">
            <b-col>
              <b-dropdown @show="refreshMatchRequests(id)" text="Requests">
                <b-dropdown-item v-show="refreshingMatchRequests">
                  <b-spinner></b-spinner>
                </b-dropdown-item>

                <b-dropdown-item
                  v-show="!refreshingMatchRequests"
                  v-for="(item, index) in matchRequests"
                  :key="index"
                  @click.prevent="doSendRequest({ request: item })"
                >{{ item }}</b-dropdown-item>
              </b-dropdown>
            </b-col>
          </b-row>

          <b-row class="pb-3">
            <b-col>
              <highlight-code
                lang="yaml"
                class="rounded"
              >
                {{ lastResponse | toYaml }}
              </highlight-code>
            </b-col>
          </b-row>

          <b-row class="pb-3">
            <b-col>
              <highlight-code lang="yaml" class="rounded">
                {{ { players: match.players } | toYaml }}
              </highlight-code>
            </b-col>
          </b-row>
          
          <b-row ref="match-deals" class="match-deals">
            <b-col>
              <highlight-code lang="yaml" class="rounded">
                {{ { deals: match.deals } | toYaml }}
              </highlight-code>
            </b-col>
          </b-row>
        </b-container>

        <b-container fluid slot="footer">
          <b-row>
            <b-col md="0">
              <b-spinner
                v-show="refreshingMatch"
              ></b-spinner>

              <b-button
                v-show="!refreshingMatch"
                @click.prevent="refreshMatch(id)"
                variant="primary"
              >
                <font-awesome-icon icon="redo"></font-awesome-icon>
              </b-button>
            </b-col>
          </b-row>
        </b-container>
      </b-card>
    </b-container>
  </div>
</template>

<script>
import Vue from 'vue'
import { mapState, mapMutations, mapActions } from 'vuex'
import VueSSE from 'vue-sse'
import BootstrapVue from 'bootstrap-vue'
import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'
import { library } from '@fortawesome/fontawesome-svg-core'
import { faRedo } from '@fortawesome/free-solid-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import VueHighlightJS from 'vue-highlight.js'
import yaml from 'highlight.js/lib/languages/yaml'
import 'highlight.js/styles/monokai.css'
import JsYaml from 'js-yaml'
import Header from '@/components/Core/Header'
import Messages from '@/components/Messages'

Vue.use(VueSSE)
Vue.use(BootstrapVue)
Vue.use(VueHighlightJS, {
  languages: { yaml }
})

library.add(faRedo)

let sseConn

export default {
  mounted() {
    this.refreshMatch(this.id)

    this.refreshMatchRequests(this.id)
    this.setLastResponse(null)

    const matchNotificationsURL = `/api/matchs/${this.id}/notifications`
    this
      .$sse(matchNotificationsURL)
      .then(sse => {
        sseConn = sse

        sse.subscribe('message', () => this.refreshMatch(this.id))
      })
      .catch(e => {
        if (e.readyState != EventSource.CLOSED)
          this.setGlobalError(`${matchNotificationsURL}: ${JSON.stringify(e)}`)
      })
  },
  updated() {
    this.$refs['match-deals'].scrollTop = this.$refs['match-deals'].scrollHeight
  },
  beforeDestroy() {
    if (sseConn)
      sseConn.close()
  },
  components: {
    FontAwesomeIcon, Header, Messages
  },
  data() {
    return {
      id: this.$route.params.id
    }
  },
  computed: {
    ...mapState([
      'refreshingMatch', 'match',
      'refreshingMatchRequests', 'matchRequests',
      'lastRequest', 'lastResponse'
    ])
  },
  methods: {
    ...mapMutations([
      'setGlobalError',
      'setLastResponse'
    ]),
    ...mapActions([
      'refreshMatch',
      'refreshMatchRequests',
      'sendRequest'
    ]),
    doSendRequest(req) {
      this.sendRequest({
        matchID: this.id,
        req
      })

      this.refreshMatch(this.id)
    }
  },
  filters: {
    toYaml(v) {
      return JsYaml.safeDump(v, {
        skipInvalid: true
      })
    }
  },
}
</script>

<style scoped>
.match-deals {
  overflow-y:scroll;
  height:40vh;
}
</style>
