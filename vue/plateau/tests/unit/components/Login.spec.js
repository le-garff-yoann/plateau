import { mount } from '@vue/test-utils'
import { BImg, BSpinner, BForm, BFormInput, BButton } from 'bootstrap-vue'
import store from '@/store'
import Login from '@/components/Login.vue'
import Messages from '@/components/Messages.vue'

describe('Login.vue', () => {
  const wrapper = mount(Login, { store })

  it('has b-img', () => expect(wrapper.find(BImg).exists()).toBe(true))
  it('has Messages', () => expect(wrapper.find(Messages).exists()).toBe(true))

  const spinner = wrapper.find(BSpinner)
  it('has a b-spinner', () => expect(spinner.exists()).toBe(true))

  const form = wrapper.find(BForm)
  it('has a b-form', () => expect(form.exists()).toBe(true))

  it('components visibility changes', () => {
    wrapper.vm.$store.commit('startLogin')
    expect(spinner.isVisible()).toBe(true)
    expect(form.isVisible()).toBe(false)

    wrapper.vm.$store.commit('stopLogin')
    expect(spinner.isVisible()).toBe(false)
    expect(form.isVisible()).toBe(true)

    wrapper.vm.$store.commit('startRegister')
    expect(spinner.isVisible()).toBe(true)
    expect(form.isVisible()).toBe(false)

    wrapper.vm.$store.commit('stopRegister')
    expect(spinner.isVisible()).toBe(false)
    expect(form.isVisible()).toBe(true)
  })

  it(
    'form has 2 b-form-input',
    () => expect(form.findAll(BFormInput).length).toBe(2)
  )

  it('form has 2 b-button', () => expect(form.findAll(BButton).length).toBe(2))
})
