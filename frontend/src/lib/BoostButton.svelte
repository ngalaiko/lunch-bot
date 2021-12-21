<script lang="ts">
  import { boosts, places, Place } from '../api'

  export let place: Place

  const onButtonClick = () => boosts.create(place.id).then(places.list).catch(alert)
</script>

<button
  type="button"
  class="text-lg p-2 px-4 rounded-md drop-shadow-md bg-blue-500 text-white"
  on:click={onButtonClick}>Boost</button
>
