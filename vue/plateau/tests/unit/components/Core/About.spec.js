import { mount } from '@vue/test-utils'
import { BMedia } from 'bootstrap-vue'
import store from '@/store'
import About from '@/components/Core/About.vue'
import Header from '@/components/Core/Header.vue'
import Messages from '@/components/Messages.vue'

describe('About.vue', () => {
  const wrapper = mount(About, { store })

  const game = {
    name: 'foo',
    description: "foo's description",
    min_players: 1,
    max_players: 2
  }
  beforeAll(() => wrapper.vm.$store.commit('setGame', game))

  it('has Header', () => expect(wrapper.find(Header).exists()).toBe(true))
  it('has Messages', () => expect(wrapper.find(Messages).exists()).toBe(true))

  const medias = wrapper.findAll(BMedia)
  it('has 4 b-media', () => expect(medias.length).toBe(4))

  it(
    "medias[0] is the game's name b-media",
    () => expect(medias.at(0).text()).toBe(`Game name ${game.name}`)
  )
  it(
    "medias[1] is the game's description b-media",
    () => expect(medias.at(1).text()).toBe(`Game description ${game.description}`)
  )
  it(
    "medias[2] is the game's min_players b-media",
    () => expect(medias.at(2).text()).toBe(`Minimum number of players ${game.min_players}`)
  )
  it(
    "medias[3] is the game's max_players b-media",
    () => expect(medias.at(3).text()).toBe(`Maximum number of players ${game.max_players}`)
  )
})
