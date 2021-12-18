import { writable } from 'svelte/store'

export type Place = {
  id: string
  name: string
  addedAt: Date
  chance: number
}

export type Roll = {
  id: string
  placeId: string
  time: Date
  place: Place
}

const store = writable({
  places: [] as Place[],
  rolls: [] as Roll[]
})

const socket = new WebSocket('ws://localhost:8000/ws')

socket.addEventListener('open', () => {
  list()
})

const storePlace = (place: Place) => {
  store.update(store => {
    store.places = store.places.filter(p => p.id !== place.id)
    store.places.push(place)
    return store
  })
}

const storeRoll = (roll: Roll) => {
  store.update(store => {
    store.rolls = store.rolls.filter(r => r.id !== roll.id)
    store.rolls.push(roll)
    return store
  })
}

const parseRoll = (roll: any): Roll => {
  return {
    id: roll.id,
    placeId: roll.placeId,
    time: new Date(roll.time),
    place: parsePlace(roll.place)
  }
}

const parsePlace = (place: any): Place => {
  return {
    id: place.id,
    name: place.name,
    addedAt: new Date(place.addedAt),
    chance: place.chance
  }
}

const list = () => {
  socket.send(
    JSON.stringify({
      id: performance.now().toString(),
      method: 'list'
    })
  )
}

export const roll = () => {
  socket.send(
    JSON.stringify({
      id: performance.now().toString(),
      method: 'roll'
    })
  )
}

export const addPlace = (name: string) => {
  socket.send(
    JSON.stringify({
      id: performance.now().toString(),
      method: 'add',
      params: {
        name
      }
    })
  )
}

socket.addEventListener('message', event => {
  const data = JSON.parse(event.data)
  if (data.places) data.places.map(parsePlace).forEach(storePlace)
  if (data.rolls) data.rolls.map(parseRoll).forEach(storeRoll)
})

export default {
  subscribe: store.subscribe
}
