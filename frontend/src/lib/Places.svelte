<script lang="ts">
  import { places } from '../api'
  import Place from './Place.svelte'
</script>

{#await places.list()}
  <p>Loading...</p>
{:then}
  <ul>
    {#each $places as place}
      <Place {place} />
    {/each}
  </ul>
{:catch e}
  <p>Error: {e.message}</p>
{/await}
