import { mount } from '@vue/test-utils'
import { BContainer, BImg, BSpinner, BFormInput, BButton } from 'bootstrap-vue'
import store from '@/store'
import Login from '@/components/Login.vue'
import Messages from '@/components/Messages.vue'

describe('Login.vue', () => {
  const wrapper = mount(Login, { store })

  it('is a b-container', () => expect(wrapper.find(BContainer).exists()).toBe(true))

  it('has b-img', () => expect(wrapper.find(BImg).exists()).toBe(true))
  it('has Messages', () => expect(wrapper.find(Messages).exists()).toBe(true))

  const spinner = wrapper.find(BSpinner)
  it('has a b-spinner', () => expect(spinner.exists()).toBe(true))

  const formInputs = wrapper.findAll(BFormInput)
  it('has 2 b-form-input', () => expect(formInputs.length).toBe(2))

  const buttons = wrapper.findAll(BButton)
  it('has 2 b-button', () => expect(buttons.length).toBe(2))
})
