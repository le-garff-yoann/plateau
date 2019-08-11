import { mount } from '@vue/test-utils'
import { BAlert } from 'bootstrap-vue'
import store from '@/store'
import Messages from '@/components/Messages.vue'

describe('Messages.vue', () => {
  const wrapper = mount(Messages, { store })

  const alerts = wrapper.findAll(BAlert)
  it('has 2 b-alert', () => expect(alerts.length).toBe(2))

  it(
    'has a hidden b-alert for messages',
    () => expect(alerts.at(0).isVisible()).toBe(false)
  )

  it(
    'has a hidden b-alert for errors',
    () => expect(alerts.at(1).isVisible()).toBe(false)
  )
  
  const globalMessage = 'msg'
  it(
    `alerts[0] is for messages and has inner text matching ${globalMessage}`,
    () => {
      wrapper.vm.setGlobalMessage(globalMessage)

      expect(alerts.at(0).isVisible()).toBe(true)
      expect(alerts.at(0).text().includes(globalMessage)).toBe(true)
    }
  )

  const globalError = 'msg'
  it(
    `alerts[1] is for errors and has inner inner text matching ${globalError}`,
    () => {
      wrapper.vm.setGlobalError(globalError)

      expect(alerts.at(1).isVisible()).toBe(true)
      expect(alerts.at(1).text().includes(globalError)).toBe(true)
    }
  )
})
