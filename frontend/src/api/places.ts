import { writable } from 'svelte/store'
import { websocket } from './protocols'

export type Place = {
  id: string
  name: string
  addedAt: Date
  chance: number
}

const store = writable<Place[]>([])

const storeResponse = (response: any) => {
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

const create = async (name: string): Promise<void> => {
  await websocket.open()
  const response = await websocket.request({ method: 'places/create', params: { name } })
  if (response.error) throw new Error(response.error)
  response.place && store.update(places => places.concat(response.place))
  storeResponse(response)
}

const list = async (): Promise<void> => {
  await websocket.open()
  const response = await websocket.request({ method: 'places/list' })
  if (response.error) throw new Error(response.error)
  storeResponse(response)
}

export default {
  list,
  create,
  subscribe: store.subscribe
}
