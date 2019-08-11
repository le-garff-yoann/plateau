import { mount } from '@vue/test-utils'
import store from '@/store'
import App from '@/App.vue'

describe('App.vue', () => {
  const wrapper = mount(App, { store })

  it('is a div', () => expect(wrapper.find('div').exists()).toBe(true))
})
