import { mount } from '@vue/test-utils'
import { BCard, BButton, BDropdown, BDropdownItem, BTab, BListGroup, BListGroupItem } from 'bootstrap-vue'
import store from '@/store'
import Match from '@/components/Core/Match.vue'
import Header from '@/components/Core/Header.vue'
import Messages from '@/components/Messages.vue'

describe('Match.vue', () => {
  Match.mounted = () => {}

  const wrapper = mount(Match, {
    store,
    stubs: [ 'Header', 'Messages' ],
    props: { id: 1 }
  })

  it('has Header', () => expect(wrapper.find(Header).exists()).toBe(true))
  it('has Messages', () => expect(wrapper.find(Messages).exists()).toBe(true))

  const card = wrapper.find(BCard)
  it('has a b-card', () => expect(card.exists()).toBe(true))

  const button = card.find(BButton)
  it('card has a b-button', () => expect(button.exists()).toBe(true))

  const dropdown = card.find(BDropdown)
  it('card has a b-dropdown', () => expect(dropdown.exists()).toBe(true))

  const matchRequests = [ 'bar', 'foo' ]
  it(`dropdown includes ${dropdown}`, () => {
    wrapper.vm.$store.commit('setMatchRequests', matchRequests)

    expect(dropdown.findAll(BDropdownItem).filter(x => !x.isVisible()).length).toBe(1)

    const dropdownItems = dropdown.findAll(BDropdownItem).filter(x => x.isVisible())
    expect(dropdownItems.length).toBe(matchRequests.length)

    for (let i = 0; i < matchRequests.length; i++) {
      expect(dropdownItems.at(i).text()).toBe(matchRequests[i])
    }
  })

  const tabs = wrapper.findAll(BTab)
  it('card has a 2 b-tab', () => expect(tabs.length).toBe(2))

  it(
    'tab[0] has a b-listgroup for deals',
    () => expect(tabs.at(0).find(BListGroup).exists()).toBe(true)
  )
  it(
    'tab[1] has a b-listgroup for players',
    () => expect(tabs.at(1).find(BListGroup).exists()).toBe(true)
  )

  const match = {
    match: { ended_at: null },
    players: [ 'foo', 'bar' ],
    deals: [ {}, {} ]
  }
  it(`tab[0] has ${match.deals.length} deals and tab[1] has ${match.players.length} players`, () => {
    wrapper.vm.$store.commit('setMatch', match)

    expect(tabs.at(0).find(BListGroup).findAll(BListGroupItem).length).
      toBe(match.deals.length)
    expect(tabs.at(1).find(BListGroup).findAll(BListGroupItem).length).
      toBe(match.players.length)
  })
})
