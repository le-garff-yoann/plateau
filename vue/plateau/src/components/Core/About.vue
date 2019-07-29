<template>
  <div>
    <Header />

    <Messages class="pt-4" />

    <b-container class="pt-4 pb-4">
      <b-row align-h="center">
        <b-col>
          <ul class="list-unstyled">
            <b-media tag="li">
              <h5 class="mt-0 mb-1">Game name</h5>
              <p class="mb-0">{{ game.name }}</p>
            </b-media>

            <b-media tag="li" class="my-4">
              <h5 class="mt-0 mb-1">Game description</h5>
              <p class="mb-0">{{ game.description }}</p>
            </b-media>

            <b-media tag="li" class="my-4">
              <h5 class="mt-0 mb-1">Minimum number of players</h5>
              <p class="mb-0">{{ game.min_players }}</p>
            </b-media>

            <b-media tag="li" class="my-4">
              <h5 class="mt-0 mb-1">Maximum number of players</h5>
              <p class="mb-0">{{ game.max_players }}</p>
            </b-media>
          </ul>
        </b-col>
      </b-row>
    </b-container>
  </div>
</template>

<script>
import Vue from 'vue'
import { mapState } from 'vuex'
import BootstrapVue from 'bootstrap-vue'
import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'
import Header from '@/components/Core/Header'
import Messages from '@/components/Messages'

Vue.use(BootstrapVue)

export default {
  components: {
    Header, Messages
  },
  computed: {
    ...mapState([
      'game'
    ])
  }
}
</script>

<style scoped>
.list-unstyled {
  text-align: left;
}
</style>
