import { writable } from 'svelte/store'
import { http } from './protocols'

export type User = {
  id: string
  name: string
}

const store = writable<User | null>(null)

const getMe = async (): Promise<void> =>
  await http
    .get('users/me')
    .then(user => {
      store.set({
        id: user.id,
        name: user.name
      })
    })
    .catch((e: Error) => {
      if (e.message.startsWith('401')) return
      throw e
    })

const slackOAuth = async (code: string, redirectUri: string): Promise<void> => {
  const user = await http.post('oauth/slack', { code, redirectUri })
  store.set({
    id: user.id,
    name: user.name
  })
}

const logout = async (): Promise<void> => {
  await http.post('users/logout')
  store.set(null)
}

export default {
  getMe,
  slackOAuth,
  logout,
  subscribe: store.subscribe
}
