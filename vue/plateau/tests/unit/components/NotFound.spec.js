import { mount } from '@vue/test-utils'
import NotFound from '@/components/NotFound.vue'

describe('NotFound.vue', () => {
  const wrapper = mount(NotFound)

  const header = wrapper.find('h6')
  const headerText = 'Oops! Page not found'

  it(`has a header with inner text set at ${headerText}`, () => {
    expect(header.exists()).toBe(true)
    expect(header.text()).toBe(headerText)
  })
})
