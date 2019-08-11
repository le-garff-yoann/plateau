import { mount } from '@vue/test-utils'
import { BButton, BSpinner, BTable, BModal } from 'bootstrap-vue'
import store from '@/store'
import Matchs from '@/components/Core/Matchs.vue'
import Header from '@/components/Core/Header.vue'
import Messages from '@/components/Messages.vue'

describe('Matchs.vue', () => {
  const wrapper = mount(Matchs, {
    store,
    stubs: [ 'Header', 'Messages' ]
  })

  it('has Header', () => expect(wrapper.find(Header).exists()).toBe(true))
  it('has Messages', () => expect(wrapper.find(Messages).exists()).toBe(true))

  const buttons = wrapper.findAll(BButton)
  it('has 2 b-buttons', () => expect(buttons.length).toBe(2))

  it(
    'has a "new match" b-button',
    () => expect(buttons.at(1).text()).toBe('New')
  )

  const spinner = wrapper.find(BSpinner)
  it('has a b-spinner', () => expect(spinner.exists()).toBe(true))

  const table = wrapper.find(BTable)
  it('has a b-table', () => expect(table.exists()).toBe(true))

  it('components visibility changes', () => {
    wrapper.vm.$store.commit('startSetMatchs')
    expect(spinner.isVisible()).toBe(true)
    expect(table.isVisible()).toBe(false)

    wrapper.vm.$store.commit('stopSetMatchs')
    expect(spinner.isVisible()).toBe(false)
    expect(table.isVisible()).toBe(true)
  })

  const id = 'foo'
  it(
    `b-table has one match identified as ${id}`,
    () => {
      wrapper.vm.$store.commit('setMatchs', [{ id }])

      expect(table.text().includes(id)).toBe(true)
    }
  )

  it(
    'has a "new match" b-modal',
    () => expect(wrapper.find(BModal).exists()).toBe(true)
  )
})
