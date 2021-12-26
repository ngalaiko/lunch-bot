<script lang="ts">
  import { createEventDispatcher } from 'svelte'
  import { Button } from '../atoms'
  import type { Place } from '../../api'

  export let place: Place

  const dispatch = createEventDispatcher()
  const onButtonClick = () => dispatch('boost', place)
</script>

<div class="flex">
  <p>
    {place.name} - {place.chance * 100}%
  </p>
  <Button size="sm" on:click={onButtonClick}>Boost</Button>
</div>
