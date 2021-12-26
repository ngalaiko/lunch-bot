<script lang="ts">
  import { places, boosts } from '../../api'
  import Place from '../molecules/Place.svelte'

  const handleOnBoost = (e: CustomEvent) =>
    boosts.create(e.detail.id).then(places.list).catch(alert)
</script>

{#await places.list()}
  <p>Loading...</p>
{:then}
  <ul>
    {#each $places as place}
      <Place on:boost={handleOnBoost} {place} />
    {/each}
  </ul>
{:catch e}
  <p>Error: {e.message}</p>
{/await}
