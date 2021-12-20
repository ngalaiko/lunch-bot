import { writable } from 'svelte/store'
import socket from './socket'

export type Place = {
  id: string
  name: string
  addedAt: Date
  chance: number
}

const store = writable<Place[]>([])

const list = async (): Promise<void> => {
  await socket.open()
  const response = await socket.sendRequest({ method: 'places/list' })
  if (response.error) throw new Error(response.error)
  response.places &&
    response.places
      .map((place: any): Place => {
        return {
          id: place.id,
          name: place.name,
          addedAt: new Date(place.addedAt),
          chance: place.chance
        }
      })
      .forEach((place: Place) =>
        store.update(places => places.filter(p => p.id !== place.id).concat(place))
      )
}

export default {
  list,
  subscribe: store.subscribe
}
