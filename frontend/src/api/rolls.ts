import { writable } from 'svelte/store'
import socket from './socket'
import type { Place } from './places'

export type Roll = {
  id: string
  placeId: string
  time: Date
  place: Place
}

const store = writable<Roll[]>([])

const storeResponse = (response: any) => {
  response.rolls &&
    response.rolls
      .map((roll: any): Roll => {
        return {
          id: roll.id,
          placeId: roll.placeId,
          time: new Date(roll.time),
          place: {
            id: roll.place.id,
            name: roll.place.name,
            addedAt: new Date(roll.place.addedAt),
            chance: roll.place.chance
          }
        }
      })
      .forEach((roll: Roll) =>
        store.update(rolls => rolls.filter(r => r.id !== roll.id).concat(roll))
      )
}

const list = async (): Promise<void> => {
  await socket.open()
  const response = await socket.sendRequest({ method: 'rolls/list' })
  if (response.error) throw new Error(response.error)
  storeResponse(response)
}

const create = async (): Promise<void> => {
  await socket.open()
  const response = await socket.sendRequest({ method: 'rolls/create' })
  if (response.error) throw new Error(response.error)
  storeResponse(response)
}

export default {
  create,
  list,
  subscribe: store.subscribe
}
