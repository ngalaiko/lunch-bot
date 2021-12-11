import { derived } from 'svelte/store'
import store from './api'

const places = derived(store, $store => {
  return $store.places
})

export default {
  subscribe: places.subscribe
}
