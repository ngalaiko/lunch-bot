<script lang="ts">
  import { createEventDispatcher } from 'svelte'
  import { Button } from '../atoms'
  import type { Place } from '../../api'

  export let place: Place

  const dispatch = createEventDispatcher()
  const onButtonClick = () => dispatch('boost', place)

  const humanChance = (chance: number) => {
    chance = chance * 100
    return `${chance.toFixed(2)}%`
  }
</script>

<span class="flex items-center justify-between">
  <p>
    {place.name}: {humanChance(place.chance)}
  </p>
  <Button size="sm" on:click={onButtonClick}>Boost</Button>
</span>
