import { mount } from '@vue/test-utils'
import { BNavbar, BNavbarBrand, BNavItem, BNavItemDropdown, BDropdownItem, CardPlugin } from 'bootstrap-vue'
import store from '@/store'
import Header from '@/components/Core/Header.vue'

describe('Header.vue', () => {
  const wrapper = mount(Header, { store })

  const navbar = wrapper.find(BNavbar)
  it('is a b-navbar', () => expect(navbar.exists()).toBe(true))

  const navbarBrand = navbar.find(BNavbarBrand)
  it(
    'navbar has a b-navbar-brand',
    () => expect(navbarBrand.exists()).toBe(true)
  )

  const gameName = 'foo'
  it(`b-navbar-brand text is set to ${gameName}`, () => {
    wrapper.vm.$store.commit('setGame', { name: gameName })

    expect(navbarBrand.text()).toBe(gameName)
  })

  const links = navbar.findAll(BNavItem)
  it('navbar has a 3 b-link', () => expect(links.length).toBe(3))

  it(
    `links[0] is the "Matchs" b-link`,
    () => expect(links.at(0).text()).toBe('Matchs')
  )
  it(
    `links[1] is the "Players" b-link`,
    () => expect(links.at(1).text()).toBe('Players')
  )
  it(
    `links[2] is the "About" b-link`,
    () => expect(links.at(2).text()).toBe('About')
  )

  const userDropdownMenu = navbar.find(BNavItemDropdown)
  it('navbar has a b-nav-item-dropdown', () => expect(userDropdownMenu.exists()).toBe(true))

  const username = 'bar'
  it(
    `b-nav-item-dropdown text matches "null"`,
    () => expect(userDropdownMenu.text().includes('null')).toBe(true)
  )
  it(`b-nav-item-dropdown text matches ${username}`, () => {
    wrapper.vm.$store.commit('setLocalCookie', { username })

    expect(userDropdownMenu.text().includes(username)).toBe(true)
  })
})
