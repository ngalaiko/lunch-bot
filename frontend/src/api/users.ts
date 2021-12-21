import { writable } from 'svelte/store'
import { rest } from './protocols'

export type User = {
  id: string
  name: string
}

const store = writable<User | null>(null)

const getMe = async (): Promise<void> => {
  const user = await rest.get('api/users/me')
  store.set({
    id: user.id,
    name: user.name
  })
}

export default {
  getMe,
  subscribe: store.subscribe
}
