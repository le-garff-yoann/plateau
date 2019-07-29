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
    `has a visible and non-empty b-alert for messages with inner text set at ${globalMessage}`,
    () => {
      wrapper.vm.setGlobalMessage(globalMessage)

      expect(alerts.at(0).isVisible()).toBe(true)
      expect(alerts.at(0).text()).toBe(`×${globalMessage}`)
    }
  )

  const globalError = 'msg'
  it(
    `has a visible and non-empty b-alert for errors with inner text set at ${globalError}`,
    () => {
      wrapper.vm.setGlobalError(globalError)

      expect(alerts.at(1).isVisible()).toBe(true)
      expect(alerts.at(1).text()).toBe(`×${globalError}`)
    }
  )
})
