import { restUri } from './api'
import { writable } from 'svelte/store'

export type User = {
  id: string
  name: string
}

const store = writable<User | null>(null)

const getMe = async (): Promise<void> => {
  const response = await fetch(`${restUri}/api/users/me`, {
    method: 'GET',
    credentials: 'include'
  })
  if (response.status !== 200) throw new Error('failed to get user')
  const body = await response.json()
  store.set({
    id: body.id,
    name: body.name
  })
}

export default {
  getMe,
  subscribe: store.subscribe
}
