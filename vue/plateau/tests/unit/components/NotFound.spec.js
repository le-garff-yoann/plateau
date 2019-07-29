import { mount } from '@vue/test-utils'
import { BContainer } from 'bootstrap-vue'
import NotFound from '@/components/NotFound.vue'

describe('NotFound.vue', () => {
  const wrapper = mount(NotFound)

  it('is a b-container', () => expect(wrapper.find(BContainer).exists()).toBe(true))

  const header = wrapper.find('h6')
  const headerText = 'Oops! Page not found'

  it(`has a header with inner text set at ${headerText}`, () => {
    expect(header.exists()).toBe(true)
    expect(header.text()).toBe(headerText)
  })
})
