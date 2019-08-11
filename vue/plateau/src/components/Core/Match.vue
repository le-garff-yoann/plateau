<template>
  <div>
    <Header />

    <Messages class="pt-4" />

    <b-container class="pt-4 pb-4" fluid>
      <b-card
        no-body
        class="text-left"
        v-bind:border-variant="match.match && match.match.ended_at == null ? 'success' : 'dark'"
      >
        <b-button-group slot="header">
          <b-button
            @click.prevent="refreshMatch(id)"
            variant="primary"
          >
            <font-awesome-icon
              icon="redo"
              v-show="!refreshingMatch"
            ></font-awesome-icon>

            <b-spinner
              small
              v-show="refreshingMatch"
            ></b-spinner>
          </b-button>

          <b-dropdown text="Requests" @show="refreshMatchRequests(id)">
            <b-dropdown-item v-show="refreshingMatchRequests">
              <b-spinner></b-spinner>
            </b-dropdown-item>

            <b-dropdown-item
              v-for="(item, index) in matchRequests.slice().sort()"
              :key="index"
              @click.prevent="doSendRequest({ request: item })"
              v-show="!refreshingMatchRequests"
            >{{ item }}</b-dropdown-item>
          </b-dropdown>
        </b-button-group>

        <b-tabs pills card end>
          <b-tab
            ref="match-deals"
            title="Deals"
            class="tab"
            no-body
          >
            <highlight-code lang="yaml">
              {{ lastResponse | toYaml }}
            </highlight-code>

            <b-list-group flush v-show="match.deals.length > 0">
              <b-list-group-item
                v-for="(item, index) in match.deals.slice().reverse()"
                :key="index"
              >
                <highlight-code
                  v-if="match.match && match.match.ended_at == null && index == 0"
                  lang="yaml"
                  class="active-deal"
                >
                  {{ item | toYaml }}
                </highlight-code>

                <highlight-code
                  v-else
                  lang="yaml"
                >
                  {{ item | toYaml }}
                </highlight-code>
              </b-list-group-item>
            </b-list-group>
          </b-tab>

          <b-tab
            title="Players"
            class="tab"
          >
            <b-list-group flush>
              <b-list-group-item
                v-for="(item, index) in match.players.slice().sort()"
                :key="index"
              >{{ item }}</b-list-group-item>
            </b-list-group>
          </b-tab>
        </b-tabs>
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

export default {
  mounted() {
    this.refreshMatch(this.id)

    this.refreshMatchRequests(this.id)
    this.setLastResponse(null)
  },
  beforeDestroy() {
    if (this.sseConn)
      this.sseConn.close()
  },
  components: {
    FontAwesomeIcon, Header, Messages
  },
  props: {
    id: String
  },
  data() {
    return {
      sseConn: null,
      sseSingleton: false
    }
  },
  computed: {
    ...mapState([
      'refreshingMatch', 'match',
      'refreshingMatchRequests', 'matchRequests',
      'lastResponse'
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
  watch: {
    refreshingMatch(val) {
      if (!this.sseSingleton && !val && this.match.match.ended_at == null) {
        this.sseSingleton = true

        const matchNotificationsURL = `/api/matchs/${this.id}/notifications`
        this
          .$sse(matchNotificationsURL)
          .then(sse => {
            this.sseConn = sse

            sse.subscribe('message', () => {
              this.refreshMatch(this.id)

              this.refreshMatchRequests(this.id)
            })
          })
          .catch(e => {
            if (e.readyState != EventSource.CLOSED)
              this.setGlobalError(`${matchNotificationsURL}: ${JSON.stringify(e)}`)
          })
      }
    }
  },
  filters: {
    toYaml(val) {
      return JsYaml.safeDump(val, {
        skipInvalid: true
      })
    }
  },
}
</script>

<style scoped>
.tab {
  overflow-y: scroll;
  height: 75vh;
}

.active-deal {
  border: 3px solid green;
}
</style>
