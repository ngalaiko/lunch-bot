import { url as apiURL } from './api'

export type User = {
  id: string
}

const getMe = async (): Promise<User | null> => {
  const response = await fetch(`${apiURL}/api/users/me`, {
    method: 'GET',
    credentials: 'include'
  })
  if (response.status !== 200) return null
  const body = await response.json()
  return {
    id: body.id
  }
}

export default {
  getMe
}
