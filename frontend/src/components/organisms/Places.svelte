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

{#await places.list()}
  <p>Loading...</p>
{:then}
  <ul>
    {#each $places as place}
      <li><Place on:boost={handleOnBoost} {place} /></li>
    {/each}
    <li><NewPlaceForm on:submit={handleOnSubmit} bind:name={newPlaceName} /></li>
  </ul>
{:catch e}
  <p>Error: {e.message}</p>
{/await}
