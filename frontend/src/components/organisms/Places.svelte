<script lang="ts">
  import { places, boosts } from '../../api'
  import { Place, NewPlaceForm } from '../molecules'

  const handleOnBoost = (e: CustomEvent) =>
    boosts.create(e.detail.id).then(places.list).catch(alert)

  let newPlaceName = ''
  const handleOnSubmit = (e: CustomEvent) => {
    if (newPlaceName.length > 0)
      places
        .create(e.detail.name)
        .then(places.list)
        .catch(alert)
        .finally(() => {
          newPlaceName = ''
        })
  }
</script>

<div class="flex flex-col items-center">
  {#await places.list()}
    <p>Loading...</p>
  {:then}
    <ul class="flex flex-col items-stretch space-y-2">
      <li><NewPlaceForm on:submit={handleOnSubmit} bind:name={newPlaceName} /></li>
      {#each $places as place}
        <li><Place on:boost={handleOnBoost} {place} /></li>
      {/each}
    </ul>
  {:catch e}
    <p>Error: {e.message}</p>
  {/await}
</div>
