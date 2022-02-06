import { writable } from 'svelte/store'
import { websocket } from './protocols'
import type { Place } from './places'
import { parseJSON as parsePlaceJSON } from './places'
import type { User } from './users'
import { parseJSON as parseUserJSON } from './users'

export type Roll = {
  __typename: 'Roll'

  placeId: string
  userId: string
  user: User
  time: Date
  place: Place
}

const sameRoll = (a: Roll, b: Roll): boolean =>
  a.placeId === b.placeId && a.userId === b.userId

const store = writable<Roll[]>([])

const parseJSON = (data: any): Roll => {
  return {
    __typename: 'Roll',
    placeId: data.placeId,
    time: new Date(data.time),
    userId: data.userId,
    user: parseUserJSON(data.user),
    place: parsePlaceJSON(data.place)
  }
}

const storeResponse = (response: any) => {
  response.rolls &&
    response.rolls.map(parseJSON).forEach((roll: Roll) =>
      store.update(rolls =>
        rolls
          .filter(r => !sameRoll(r, roll))
          .concat(roll)
          .sort((a, b) => b.time.getTime() - a.time.getTime())
      )
    )
}

let listPromise: Promise<any> | undefined
const list = async (): Promise<void> => {
  await websocket.open()
  if (listPromise) return await listPromise
  listPromise = websocket.request({ method: 'rolls/list' })
  const response = await listPromise
  if (response.error) throw new Error(response.error)
  storeResponse(response)
}

const create = async (): Promise<void> => {
  await websocket.open()
  const response = await websocket.request({ method: 'rolls/create' })
  if (response.error) throw new Error(response.error)
  storeResponse(response)
}

websocket.onMessage(storeResponse)

export default {
  create,
  list,
  subscribe: store.subscribe
}
