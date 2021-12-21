import { writable } from 'svelte/store'
import socket from './socket'

export type Boost = {
  id: string
  userId: string
  placeId: string
  time: Date
}

const store = writable<Boost[]>([])

const storeResponse = (response: any) => {
  response.boosts &&
    response.boosts
      .map((boost: any) => {
        return {
          id: boost.id,
          userId: boost.userId,
          placeId: boost.placeId,
          time: new Date(boost.time)
        }
      })
      .forEach((boost: Boost) => {
        store.update(boosts => boosts.filter(b => b.id !== boost.id).concat(boost))
      })
}

const create = async (placeId: string): Promise<void> => {
  await socket.open()
  const response = await socket.sendRequest({
    method: 'boosts/create',
    params: { placeId }
  })
  if (response.error) throw new Error(response.error)
  storeResponse(response)
}

export default {
  create,
  subscribe: store.subscribe
}
