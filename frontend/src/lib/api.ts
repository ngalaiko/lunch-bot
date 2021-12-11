import { writable } from 'svelte/store'

export type Place = {
  id: string
  name: string
  addedAt: Date
  chance: number
}

const store = writable({
  places: [] as Place[]
})

const socket = new WebSocket('ws://localhost:8000/ws')

socket.addEventListener('open', function () {
  listPlaces()
})

const storePlace = (place: Place) => {
  store.update(store => {
    store.places = store.places.filter(p => p.id !== place.id)
    store.places.push(place)
    return store
  })
}

const parsePlace = (place: any): Place => {
  return {
    id: place.id,
    name: place.name,
    addedAt: new Date(place.addedAt),
    chance: place.chance
  }
}

const listPlaces = () => {
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

socket.addEventListener('message', function (event) {
  const data = JSON.parse(event.data)
  console.log(data)
  data.places.map(parsePlace).forEach(storePlace)
})

export default {
  subscribe: store.subscribe
}
