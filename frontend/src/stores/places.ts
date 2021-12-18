import { derived } from 'svelte/store'
import store from '../api/protocol'

const places = derived(store, $store => {
  return $store.places
    .sort((a, b) => {
      if (a.name < b.name) return -1
      if (a.name > b.name) return 1
      return 0
    })
    .sort((a, b) => b.chance - a.chance)
})

export default {
  subscribe: places.subscribe
}
