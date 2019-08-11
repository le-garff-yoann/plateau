import { mount } from '@vue/test-utils'
import { BButton, BSpinner, BTable } from 'bootstrap-vue'
import store from '@/store'
import Players from '@/components/Core/Players.vue'
import Header from '@/components/Core/Header.vue'
import Messages from '@/components/Messages.vue'

describe('Players.vue', () => {
  const wrapper = mount(Players, { store })

  it('has Header', () => expect(wrapper.find(Header).exists()).toBe(true))
  it('has Messages', () => expect(wrapper.find(Messages).exists()).toBe(true))

  it(
    'has a refresh b-button',
    () => expect(wrapper.find(BButton).exists()).toBe(true)
  )

  const spinner = wrapper.find(BSpinner)
  it('has a b-spinner', () => expect(spinner.exists()).toBe(true))

  const table = wrapper.find(BTable)
  it('has a b-table', () => expect(table.exists()).toBe(true))

  it('components visibility changes', () => {
    wrapper.vm.$store.commit('startSetPlayers')
    expect(spinner.isVisible()).toBe(true)
    expect(table.isVisible()).toBe(false)

    wrapper.vm.$store.commit('stopSetPlayers')
    expect(spinner.isVisible()).toBe(false)
    expect(table.isVisible()).toBe(true)
  })

  const name = 'foo'
  it(
    `b-table has one player named ${name}`,
    () => {
      wrapper.vm.$store.commit('setPlayers', [{ name }])

      expect(table.text().includes(name)).toBe(true)
    }
  )
})
