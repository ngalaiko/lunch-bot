<script lang="ts">
  import { places } from '../api'
  import Place from './Place.svelte'
</script>

<ul>
  {#await places.list()}
    Loading...
  {:then}
    {#each $places as place}
      <Place {place} />
    {/each}
  {:catch e}
    <p>Error: {e.message}</p>
  {/await}
</ul>
