import { writable } from 'svelte/store'
import { websocket } from './protocols'
import type { User } from './users'
import { parseJSON as parseUserJSON } from './users'

export type Place = {
  __typename: 'Place'

  id: string
  name: string
  chance: number
  time: Date
  userId?: string
  user?: User
}

const store = writable<Place[]>([])

export const parseJSON = (data: any): Place => {
  return {
    __typename: 'Place',

    id: data.id,
    name: data.name,
    chance: data.chance,
    time: new Date(data.time),
    userId: data.userId,
    user: data.user ? parseUserJSON(data.user) : undefined
  }
}

const storeResponse = (response: any) => {
  response.places &&
    response.places
      .map(parseJSON)
      .forEach((place: Place) =>
        store.update(places => places.filter(p => p.id !== place.id).concat(place))
      )
}

const create = async (name: string): Promise<void> => {
  await websocket.open()
  const response = await websocket.request({ method: 'places/create', params: { name } })
  if (response.error) throw new Error(response.error)
  response.place && store.update(places => places.concat(response.place))
  storeResponse(response)
}

let listPromise: Promise<any> | undefined
const list = async (): Promise<void> => {
  await websocket.open()
  if (listPromise) return await listPromise
  listPromise = websocket.request({ method: 'places/list' })
  const response = await listPromise
  if (response.error) throw new Error(response.error)
  storeResponse(response)
}

websocket.onMessage(storeResponse)

export default {
  list,
  create,
  subscribe: store.subscribe
}
