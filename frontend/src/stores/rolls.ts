import { derived } from 'svelte/store'
import store from '../api/protocol'

const rolls = derived(store, $store => {
  return $store.rolls.sort((a, b) => {
    return a.time > b.time ? -1 : 1
  })
})

export default {
  subscribe: rolls.subscribe
}
