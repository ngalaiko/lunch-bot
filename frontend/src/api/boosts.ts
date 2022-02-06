import { writable } from 'svelte/store'
import { websocket } from './protocols'
import type { Place } from './places'
import { parseJSON as parsePlaceJSON } from './places'
import type { User } from './users'
import { parseJSON as parseUserJSON } from './users'

export type Boost = {
  __typename: 'Boost'

  time: Date
  userId: string
  user: User
  placeId: string
  place: Place
}

const isSameBoost = (boost: Boost, other: Boost): boolean => {
  return (
    boost.time.getTime() === other.time.getTime() &&
    boost.placeId === other.placeId &&
    boost.userId === other.userId
  )
}

const store = writable<Boost[]>([])

const parseJSON = (data: any): Boost => {
  return {
    __typename: 'Boost',

    time: new Date(data.time),
    placeId: data.placeId,
    place: parsePlaceJSON(data.place),
    userId: data.userId,
    user: parseUserJSON(data.user)
  }
}

const storeResponse = (response: any) => {
  response.boosts &&
    response.boosts.map(parseJSON).forEach((boost: Boost) => {
      store.update(boosts => boosts.filter(b => !isSameBoost(b, boost)).concat(boost))
    })
}

const create = async (placeId: string): Promise<void> => {
  await websocket.open()
  const response = await websocket.request({
    method: 'boosts/create',
    params: { placeId }
  })
  if (response.error) throw new Error(response.error)
  storeResponse(response)
}

let listPromise: Promise<any> | undefined
const list = async (): Promise<void> => {
  await websocket.open()
  if (listPromise) return await listPromise
  listPromise = await websocket.request({ method: 'boosts/list' })
  const response = await listPromise
  if (response.error) throw new Error(response.error)
  storeResponse(response)
}

websocket.onMessage(storeResponse)

export default {
  create,
  list,
  subscribe: store.subscribe
}
