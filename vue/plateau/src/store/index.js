import Vue from 'vue'
import Vuex from 'vuex'
import VueCookies from 'vue-cookies'
import axios from 'axios'
import consts from '@/constants/index.js'

Vue.use(Vuex)
Vue.use(VueCookies)

export default new Vuex.Store({
  state: {
    sessionCookie: null,
    globalMessage: null,
    globalError: null,
    loggingIn: false,
    registerIn: false,
    game: {},
    refreshingMatchs: false,
    matchs: [],
    refreshingPlayers: false,
    players: [],
    refreshingMatch: false,
    match: {},
    refreshingMatchRequests: false,
    matchRequests: [],
    lastResponse: null
  },
  mutations: {
    setSessionCookie(state, val) { state.sessionCookie = val },
    setGlobalMessage(state, val) { state.globalMessage = val },
    setGlobalError(state, err) { state.globalError = err },
    startLogin(state) { state.loggingIn = true },
    stopLogin(state) { state.loggingIn = false },
    startRegister(state) { state.registerIn = true },
    stopRegister(state) { state.registerIn = false },
    setGame(state, val) { state.game = val },
    startSetMatchs(state) { state.refreshingMatchs = true },
    stopSetMatchs(state) { state.refreshingMatchs = false },
    setMatchs(state, val) { state.matchs = val },
    startSetPlayers(state) { state.refreshingPlayers = true },
    stopSetPlayers(state) { state.refreshingPlayers = false },
    setPlayers(state, val) { state.players = val },
    startSetMatch(state) { state.refreshingMatch = true },
    stopSetMatch(state) { state.refreshingMatch = false },
    setMatch(state, val) { state.match = val },
    startSetMatchRequests(state) { state.refreshingMatchRequests = true },
    stopSetMatchRequests(state) { state.refreshingMatchRequests = false },
    setMatchRequests(state, val) { state.matchRequests = val },
    setLastResponse(state, val) { state.lastResponse = val }
  },
  actions: {
    login({ commit }, userinfo) {
      commit('startLogin')

      axios
        .post('/user/login', userinfo)
        .then(() => {
          commit('setSessionCookie', window.$cookies.get(consts.serverName))

          window.$cookies.set(consts.localCookieName, { username: userinfo.username })
        })
        .catch(err => commit('setGlobalError', err))
        .then(() => commit('stopLogin'))
    },
    register({ commit }, userinfo) {
      commit('startRegister')

      axios
        .post('/user/register', userinfo)
        .then(() => commit('setGlobalMessage', 'Successfully registered.'))
        .catch(err => commit('setGlobalError', err))
        .then(() => commit('stopRegister'))
    },
    logout({ commit }) {
      axios
        .delete('/user/logout')
        .then(() => {
          window.$cookies.remove(consts.serverName)
          commit('setSessionCookie', window.$cookies.get(consts.serverName))

          window.$cookies.remove(consts.localCookieName)
        })
        .catch(err => commit('setGlobalError', err))
    },
    refreshGame({ commit }) {
      axios
        .get('/api/game')
        .then(res => commit('setGame', res.data))
        .catch(err => commit('setGlobalError', err))
    },
    refreshMatchs({ commit }) {
      commit('startSetMatchs')

      axios
        .get('/api/matchs')
        .then(res => {
          axios
          .all(res.data.map(id => axios.get(`/api/matchs/${id}`)))
          .then(res => commit('setMatchs', res.map(x => x.data)))
          .catch(err => commit('setGlobalError', err))
        })
        .catch(err => commit('setGlobalError', err))
        .then(() => commit('stopSetMatchs'))
    },
    createMatch({ commit, dispatch }, numberOfPlayersRequired) {
      axios
        .post('/api/matchs', { number_of_players_required: Number(numberOfPlayersRequired) })
        .then(() => dispatch('refreshMatchs'))
        .catch(err => commit('setGlobalError', err))
    },
    refreshPlayers({ commit }) {
      commit('startSetPlayers')

      axios
        .get('/api/players')
        .then(res => {
          axios
          .all(res.data.map(name => axios.get(`/api/players/${name}`)))
          .then(res => commit('setPlayers', res.map(x => x.data)))
          .catch(err => commit('setGlobalError', err))
        })
        .catch(err => commit('setGlobalError', err))
        .then(() => commit('stopSetPlayers'))
    },
    refreshMatch({ commit }, matchID) {
      commit('startSetMatch')

      axios.all([
        axios.get(`/api/matchs/${matchID}`),
        axios.get(`/api/matchs/${matchID}/players`),
        axios.get(`/api/matchs/${matchID}/deals`)
      ])
      .then(res => {
        commit('setMatch', {
          match: res[0].data,
          players: res[1].data,
          deals: res[2].data
        })
      })
      .catch(err => commit('setGlobalError', err))
      .then(() => commit('stopSetMatch'))
    },
    refreshMatchRequests({ commit }, matchID) {
      commit('startSetMatchRequests')

      axios
        .patch(`/api/matchs/${matchID}`, { request: '?'})
        .then(res => commit('setMatchRequests', res.data.body))
        .catch(err => commit('setGlobalError', err))
        .then(() => commit('stopSetMatchRequests'))
    },
    sendRequest({ commit }, data) {
      axios
        .patch(`/api/matchs/${data.matchID}`, data.req)
        .then(res => commit('setLastResponse', res.data))
        .catch(err => commit('setGlobalError', err))
    }
  }
})
