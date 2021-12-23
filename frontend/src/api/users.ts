import { writable } from 'svelte/store'
import { http } from './protocols'

export type User = {
  id: string
  name: string
}

const store = writable<User | null>(null)

const getMe = async (): Promise<void> => {
  const user = await http.get('users/me')
  store.set({
    id: user.id,
    name: user.name
  })
}

export default {
  getMe,
  subscribe: store.subscribe
}
